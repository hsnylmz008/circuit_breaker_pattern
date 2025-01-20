FROM golang:1.21-alpine AS builder

WORKDIR /app

# Gerekli paketleri yükle
RUN apk add --no-cache git

# Go modüllerini kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final image
FROM alpine:latest

WORKDIR /app

# SSL sertifikaları
RUN apk --no-cache add ca-certificates

# Builder'dan binary'yi kopyala
COPY --from=builder /app/api .

# Çalıştırma
CMD ["./api"] 