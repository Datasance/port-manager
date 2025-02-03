FROM golang:1.23-alpine3.20 AS backend

WORKDIR /port-manager

COPY ./go.* ./
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Makefile ./
COPY ./vendor ./vendor


RUN apk add --update --no-cache bash curl git make

RUN make build
RUN cp ./bin/port-manager /bin

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
WORKDIR /

RUN microdnf install shadow-utils && \
    microdnf clean all
RUN useradd --uid 10000 runner
COPY LICENSE /licenses/LICENSE

COPY --from=backend /bin/port-manager /bin
USER 10000
LABEL org.opencontainers.image.description=Port-Manager
LABEL org.opencontainers.image.source=https://github.com/datasance/port-manager
LABEL org.opencontainers.image.licenses=EPL2.0
ENTRYPOINT ["/bin/port-manager"]
