FROM --platform=$BUILDPLATFORM golang:1.21-alpine3.19 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x

COPY *.go ./
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /go/bin/pushgateway-cleaner .

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /go/bin/pushgateway-cleaner /usr/local/bin/pushgateway-cleaner

ENTRYPOINT ["pushgateway-cleaner"]
