FROM cargo.caicloud.xyz/library/golang:1.12-security as build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/db-manager
RUN if [ "$(uname -m)" = "ppc64le" ] || [ "$(uname -m)" = "aarch64" ]; then \
        apk --update add git gcc musl-dev && \
        go build -o katib-db-manager ./v1alpha3; \
    else \
        go build -o katib-db-manager ./v1alpha3; \
    fi

FROM cargo.caicloud.xyz/library/debian:stretch
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/db-manager/katib-db-manager /app/
ENTRYPOINT ["./katib-db-manager"]
CMD ["-w", "kubernetes"]
