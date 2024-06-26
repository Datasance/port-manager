/*
 *  *******************************************************************************
 *  * Copyright (c) 2023 Datasance Teknoloji A.S.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package manager

import (
	"time"

	ioclient "github.com/datasance/iofog-go-sdk/v3/pkg/client"
)

type portMap map[int]ioclient.PublicPort // Indexed by port

var pkg struct {
	controllerServiceName string
	controllerPort        int
	managerName           string
	pollInterval          time.Duration
}

func init() {
	pkg.controllerServiceName = "controller"
	pkg.controllerPort = 51121
	pkg.managerName = "port-manager"
	pkg.pollInterval = time.Second * 10
}
