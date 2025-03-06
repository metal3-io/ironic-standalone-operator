# Support FROM override
ARG BUILD_IMAGE=docker.io/golang:1.23.7@sha256:1acb493b9f9dfdfe705042ce09e8ded908ce4fb342405ecf3ca61ce7f3b168c7
ARG BASE_IMAGE=gcr.io/distroless/static:nonroot@sha256:9ecc53c269509f63c69a266168e4a687c7eb8c0cfd753bd8bfcaa4f58a90876f

# Build the manager binary
FROM $BUILD_IMAGE AS builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download -x

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/controller/ internal/controller/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
FROM $BASE_IMAGE
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]

LABEL io.k8s.display-name="Metal3 Ironic Operator" \
      io.k8s.description="Operator managing an Ironic deployment for Metal3"
