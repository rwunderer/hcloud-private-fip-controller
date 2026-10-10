#--------
# builder
#--------
FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

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
FROM gcr.io/distroless/static@sha256:58133991db06659feaabe0f4e97a35cebf15ef4ea08f8a4c6d2ee5f75e4aa6a0 AS hcloud-private-fip-controller

LABEL version=0.2.88

WORKDIR /app

USER 10000

COPY --from=builder /app/private-fip-controller /app/private-fip-controller

ENTRYPOINT ["./private-fip-controller"]
