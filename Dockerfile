# Estágio de Compilação
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o pdf-merger ./cmd/web/

# Estágio Final
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/pdf-merger .
COPY --from=builder /app/index.html .
EXPOSE 8080
CMD ["./pdf-merger"]