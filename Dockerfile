# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/server /usr/local/bin/server
COPY wait-for-contracts.sh /usr/local/bin/wait-for-contracts.sh
RUN chmod +x /usr/local/bin/wait-for-contracts.sh
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/wait-for-contracts.sh"]
CMD ["server"]
