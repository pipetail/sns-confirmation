FROM golang:1.20 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build -o main cmd/main.go

# https://edu.chainguard.dev/chainguard/chainguard-images/reference/busybox/image_specs/
FROM cgr.dev/chainguard/busybox:latest
COPY --from=builder /src/main /main
CMD ["/main"]