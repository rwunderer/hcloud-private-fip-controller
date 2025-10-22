#--------
# builder
#--------
FROM golang:1.25.3-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS builder

ARG TARGETPLATFORM
ARG TARGETOS="linux"
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT#"v"} go build -a -installsuffix cgo -ldflags '-extldflags "-static"' ./cmd/private-fip-controller

#--------
# container
#--------
FROM gcr.io/distroless/static@sha256:87bce11be0af225e4ca761c40babb06d6d559f5767fbf7dc3c47f0f1a466b92c AS hcloud-private-fip-controller

LABEL version=0.2.57

WORKDIR /app

USER 10000

COPY --from=builder /app/private-fip-controller /app/private-fip-controller

ENTRYPOINT ["./private-fip-controller"]
