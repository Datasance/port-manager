FROM golang:1.19-alpine as backend

WORKDIR /port-manager

COPY ./go.* ./
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Makefile ./
COPY ./vendor ./vendor


RUN apk add --update --no-cache bash curl git make

RUN make build
RUN cp ./bin/port-manager /bin

FROM alpine:3.19
COPY --from=backend /bin /bin
LABEL org.opencontainers.image.description Port-Manager
LABEL org.opencontainers.image.source=https://github.com/datasance/port-manager
LABEL org.opencontainers.image.licenses=EPL2.0
ENTRYPOINT ["/bin/port-manager"]
