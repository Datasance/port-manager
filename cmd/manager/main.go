package main

import (
	"encoding/json"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/datasance/port-manager/v3/internal/manager"
)

var log = zap.New()

const (
	authURLEnv                 = "KC_URL"
	realmEnv                   = "KC_REALM"
	clientIDEnv                = "KC_CLIENT"
	clientSecretEnv            = "KC_CLIENT_SECRET"
	proxyImageEnv              = "PROXY_IMAGE"
	imagePullSecretEnv         = "PULL_SECRET_NAME"
	httpProxyAddressEnv        = "HTTP_PROXY_ADDRESS"
	tcpProxyAddressEnv         = "TCP_PROXY_ADDRESS"
	routerAddressEnv           = "ROUTER_ADDRESS"
	proxyServiceAnnotationsEnv = "PROXY_SERVICE_ANNOTATIONS"
	controllerSchemeEnv        = "CONTROLLER_SCHEME"
)

type env struct {
	optional bool
	key      string
	value    string
}

func generateManagerOptions(namespace string, cfg *rest.Config) (opts []manager.Options) {
	envs := map[string]env{
		authURLEnv:                 {key: authURLEnv},
		realmEnv:                   {key: realmEnv},
		clientIDEnv:                {key: clientIDEnv},
		clientSecretEnv:            {key: clientSecretEnv},
		routerAddressEnv:           {key: routerAddressEnv},
		proxyImageEnv:              {key: proxyImageEnv},
		imagePullSecretEnv:         {key: imagePullSecretEnv, optional: true},
		httpProxyAddressEnv:        {key: httpProxyAddressEnv, optional: true},
		tcpProxyAddressEnv:         {key: tcpProxyAddressEnv, optional: true},
		proxyServiceAnnotationsEnv: {key: proxyServiceAnnotationsEnv, optional: true},
		controllerSchemeEnv:        {key: controllerSchemeEnv},
	}
	// Read env vars
	for _, env := range envs {
		env.value = os.Getenv(env.key)
		if env.value == "" && !env.optional {
			log.Error(nil, env.key+" env var not set")
			os.Exit(1)
		}
		// Store result for later
		envs[env.key] = env
	}

	opt := manager.Options{
		Namespace:               namespace,
		AuthURL:                 envs[authURLEnv].value,
		Realm:                   envs[realmEnv].value,
		ClientID:                envs[clientIDEnv].value,
		ClientSecret:            envs[clientSecretEnv].value,
		ProxyImage:              envs[proxyImageEnv].value,
		ImagePullSecret:         envs[imagePullSecretEnv].value,
		ProxyServiceType:        "LoadBalancer",
		ProxyServiceAnnotations: make(map[string]string),
		ProxyExternalAddress:    "",
		ProtocolFilter:          "",
		ProxyName:               "http-proxy", // TODO: Fix this default, e.g. iofogctl tests get svc name
		RouterAddress:           envs[routerAddressEnv].value,
		ControllerScheme:        envs[controllerSchemeEnv].value,
		Config:                  cfg,
	}

	// Set proxyServiceAnnotations if present
	if annotations := envs[proxyServiceAnnotationsEnv].value; annotations != "" {
		var annotationsMap map[string]string
		err := json.Unmarshal([]byte(annotations), &annotationsMap)
		if err != nil {
			log.Error(err, "Failed to unmarshal proxy service annotations")
			os.Exit(1)
		}
		opt.ProxyServiceAnnotations = annotationsMap
	}

	opts = append(opts, opt)
	if envs[httpProxyAddressEnv].value != "" && envs[tcpProxyAddressEnv].value != "" {
		// Update first opt
		opts[0].ProxyServiceType = "ClusterIP"
		opts[0].ProtocolFilter = "http"
		opts[0].ProxyName = "http-proxy"
		opts[0].ProxyExternalAddress = envs[httpProxyAddressEnv].value
		// Create second opt
		opt.ProxyServiceType = "ClusterIP"
		opt.ProtocolFilter = "tcp"
		opt.ProxyName = "tcp-proxy"
		opt.ProxyExternalAddress = envs[tcpProxyAddressEnv].value
		opts = append(opts, opt)
	}
	return opts
}

func generateManagers(namespace string, cfg *rest.Config) (mgrs []*manager.Manager) {
	opts := generateManagerOptions(namespace, cfg)
	// No external address provided, Manager will create Proxy LoadBalancer and single Deployment
	for idx := range opts {
		opt := &opts[idx]
		mgr, err := manager.New(opt)
		handleErr(err, "")
		mgrs = append(mgrs, mgr)
	}
	return
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Error(err, msg)
		os.Exit(1)
	}
}

// getWatchNamespace returns the Namespace the operator should be watching for changes
func getWatchNamespace() (ns string) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	ns, _ = os.LookupEnv("WATCH_NAMESPACE")
	return
}

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	handleErr(err, "")

	// Instantiate Manager(s)
	mgrs := generateManagers(getWatchNamespace(), cfg)

	// Run Managers
	for _, mgr := range mgrs {
		go mgr.Run()
	}

	// Set ready
	readyPath := "/tmp/operator-sdk-ready"
	if _, err := os.Stat(readyPath); os.IsNotExist(err) {
		file, err := os.Create(readyPath)
		handleErr(err, "Failed to create ready file")
		defer file.Close()
	}

	// Wait forever
	select {}
}
