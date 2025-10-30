# Support FROM override
ARG BUILD_IMAGE=docker.io/golang:1.24.9@sha256:5034fa44b36163a4109b71ed75c67dbdcb52c3cd9750953befe00054315d9fd2
ARG BASE_IMAGE=gcr.io/distroless/static:nonroot@sha256:9ecc53c269509f63c69a266168e4a687c7eb8c0cfd753bd8bfcaa4f58a90876f

# Build the manager binary
FROM $BUILD_IMAGE AS builder

WORKDIR /workspace
ARG LDFLAGS=-s -w -extldflags=-static
# Copy the Go Modules manifests
COPY go.mod go.sum ./
COPY api/go.mod api/go.sum api/
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download -x

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 go build -a -ldflags "${LDFLAGS}" \
    -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
FROM $BASE_IMAGE

# image.version is set during image build by automation
LABEL org.opencontainers.image.authors="metal3-dev@googlegroups.com"
LABEL org.opencontainers.image.description="Operator managing an Ironic deployment for Metal3"
LABEL org.opencontainers.image.documentation="https://book.metal3.io/irso/introduction"
LABEL org.opencontainers.image.licenses="Apache License 2.0"
LABEL org.opencontainers.image.title="Ironic Standalone Operator"
LABEL org.opencontainers.image.url="https://github.com/metal3-io/ironic-standalone-operator"
LABEL org.opencontainers.image.vendor="Metal3-io"

WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
