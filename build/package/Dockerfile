#--------
# builder
#--------
FROM golang:1.26.5-alpine@sha256:99e12cfb19b753915f9b9fdc5a99f1869a24a69d3a0955832d5702e7fa68f1be AS builder

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
FROM gcr.io/distroless/static@sha256:47b2d72ff90843eb8a768b5c2f89b40741843b639d065b9b937b07cd59b479c6 AS hcloud-private-fip-controller

LABEL version=0.2.80

WORKDIR /app

USER 10000

COPY --from=builder /app/private-fip-controller /app/private-fip-controller

ENTRYPOINT ["./private-fip-controller"]
