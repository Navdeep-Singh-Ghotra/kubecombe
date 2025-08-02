FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /product-service

FROM alpine:latest
COPY --from=builder /product-service /product-service
EXPOSE 8080 50051
ENTRYPOINT ["/product-service"]