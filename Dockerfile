FROM golang:1.24.0 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go get github.com/navdeep-singh-ghotra/kubecombe/pb
RUN CGO_ENABLED=0 GOOS=linux go build -o /kubecombe

FROM alpine:latest
COPY --from=builder /kubecombe /kubecombe
EXPOSE 8080 50051
ENTRYPOINT ["/kubecombe"]