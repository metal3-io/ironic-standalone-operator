# Support FROM override
ARG BUILD_IMAGE=docker.io/golang:1.20.11@sha256:0f3cc978eb0fe87fec189c4dc74f456cd125b43cbf94fcd645e3f4bb59bcc316
ARG BASE_IMAGE=gcr.io/distroless/static:nonroot@sha256:91ca4720011393f4d4cab3a01fa5814ee2714b7d40e6c74f2505f74168398ca9

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
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
FROM $BASE_IMAGE
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]

LABEL io.k8s.display-name="Metal3 Ironic Operator" \
      io.k8s.description="Operator managing an Ironic deployment for Metal3"
