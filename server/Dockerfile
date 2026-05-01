FROM golang:1.26 AS builder
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . .
WORKDIR /workspace/app
RUN make build

FROM scratch
COPY --from=builder /workspace/app/reflecitve-target .
USER 65534:65534
ENTRYPOINT ["/reflecitve-target"]
