#--------
# builder
#--------
FROM golang:1.25.7-alpine@sha256:81d49e1de26fa223b9ae0b4d5a4065ff8176a7d80aa5ef0bd9f2eee430afe4d7 AS builder

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
FROM gcr.io/distroless/static@sha256:972618ca78034aaddc55864342014a96b85108c607372f7cbd0dbd1361f1d841 AS hcloud-private-fip-controller

LABEL version=0.2.72

WORKDIR /app

USER 10000

COPY --from=builder /app/private-fip-controller /app/private-fip-controller

ENTRYPOINT ["./private-fip-controller"]
