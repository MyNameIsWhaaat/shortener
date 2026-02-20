FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /shortener ./cmd/shortener

# run stage
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /shortener /shortener

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/shortener"]