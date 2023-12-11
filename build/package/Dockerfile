#--------
# builder
#--------
FROM golang:alpine@sha256:feceecc0e1d73d085040a8844de11a2858ba4a0c58c16a672f1736daecc2a4ff as builder

WORKDIR /app

COPY go.mod go.sum /app/

#RUN apk add git \
#  && go mod download
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' ./cmd/private-fip-controller

#--------
# container
#--------
FROM gcr.io/distroless/static@sha256:6706c73aae2afaa8201d63cc3dda48753c09bcd6c300762251065c0f7e602b25

WORKDIR /app

USER 10000

COPY --from=builder /app/private-fip-controller /app/private-fip-controller

ENTRYPOINT ["./private-fip-controller"]
