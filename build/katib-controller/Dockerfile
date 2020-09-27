# Build the manager binary
FROM cargo.caicloud.xyz/library/golang:1.12-security as build-env

# Copy in the go src
ADD . /root/go/src/github.com/kubeflow/katib

WORKDIR /root/go/src/github.com/kubeflow/katib/cmd/katib-controller
# Build
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -a -o katib-controller  ./v1alpha3; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
        CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o katib-controller  ./v1alpha3; \
    else \
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o katib-controller  ./v1alpha3; \
    fi
# Copy the controller-manager into a thin image
FROM cargo.caicloud.xyz/library/debian:stretch
WORKDIR /app

COPY --from=build-env /root/go/src/github.com/kubeflow/katib/cmd/katib-controller/katib-controller .
ENTRYPOINT ["./katib-controller"]
