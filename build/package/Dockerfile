#--------
# builder
#--------
FROM golang:1.27.4-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS builder

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
FROM gcr.io/distroless/static@sha256:4c2cd1dee442c9a99924a378c7c126fac8fef14771579ef567c3f2bdd7a5122a AS hcloud-private-fip-controller

LABEL version=0.2.92

WORKDIR /app

USER 10000

COPY --from=builder /app/private-fip-controller /app/private-fip-controller

ENTRYPOINT ["./private-fip-controller"]
