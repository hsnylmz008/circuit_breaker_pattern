# Go Circuit Breaker Example

Bu proje, Go dilinde Circuit Breaker pattern'inin nasıl uygulanacağını gösteren bir örnek uygulamadır. Fiber web framework'ü ve Sony'nin gobreaker kütüphanesi kullanılarak geliştirilmiştir.

## Proje Yapısı

```
.
├── cmd
│   └── api
│       └── main.go           # Ana uygulama giriş noktası
├── internal
│   ├── config
│   │   └── config.go         # Konfigürasyon yönetimi
│   ├── handler
│   │   └── handler.go        # HTTP handler'ları
│   ├── router
│   │   └── router.go         # Route tanımlamaları
│   └── service
│       └── service.go        # İş mantığı katmanı
├── pkg
│   └── circuitbreaker
│       └── breaker.go        # Circuit breaker implementasyonu
├── config.json               # Konfigürasyon dosyası
├── Dockerfile               # Container yapılandırması
├── docker-compose.yml       # Docker compose yapılandırması
├── go.mod
└── README.md
```

## Circuit Breaker Pattern Nedir?

Circuit Breaker (Devre Kesici) pattern'i, dağıtık sistemlerde hata toleransını artırmak için kullanılan bir tasarım kalıbıdır. Temel amacı:

- Başarısız olma olasılığı yüksek işlemleri izole etmek
- Sistemin geri kalanını korumak
- Hızlı başarısızlık (fail-fast) sağlamak
- Sistemin kendini onarmasına olanak tanımak

### Circuit Breaker Durumları

1. **Kapalı (Closed)**: Normal çalışma durumu
   - Tüm istekler normal şekilde işlenir
   - Hata oranı izlenir

2. **Açık (Open)**: Hata durumu
   - İstekler engellenir
   - Timeout süresi sonunda yarı-açık duruma geçer

3. **Yarı-Açık (Half-Open)**: Test durumu
   - Sınırlı sayıda istek kabul edilir
   - Başarılı istekler durumu kapalıya çevirir
   - Başarısız istekler durumu açığa çevirir

## Kurulum ve Çalıştırma

### Gereksinimler
- Go 1.21 veya üstü
- Docker ve Docker Compose (opsiyonel)

### Doğrudan Çalıştırma
```bash
# Bağımlılıkları yükle
go mod download

# Uygulamayı çalıştır
go run cmd/api/main.go
```

### Docker ile Çalıştırma
```bash
# Docker compose ile çalıştır
docker-compose up --build
```

## API Endpoints

### GET /api
Test endpoint'i (simüle edilmiş dış servis çağrısı)

**Başarılı Yanıt:**
```json
{
    "result": "Success!",
    "state": "closed"
}
```

**Hata Yanıtı:**
```json
{
    "error": "service error",
    "state": "open"
}
```

### GET /health
Sağlık kontrolü endpoint'i

```json
{
    "status": "healthy",
    "circuitState": "closed"
}
```

## Konfigürasyon

`config.json` dosyası üzerinden yapılandırma:

```json
{
    "server": {
        "port": 3000,
        "host": "0.0.0.0"
    },
    "circuitBreaker": {
        "name": "API-ServiceBreaker",
        "maxRequests": 3,
        "interval": 10,
        "timeout": 30,
        "failureRatio": 0.6,
        "minimumRequests": 10
    }
}
```

### Konfigürasyon Parametreleri

#### Server
- `port`: Sunucunun çalışacağı port
- `host`: Sunucunun dinleyeceği host adresi

#### Circuit Breaker
- `maxRequests`: Yarı-açık durumda izin verilen maksimum istek sayısı
- `interval`: İstatistiklerin sıfırlanma süresi (saniye)
- `timeout`: Açık durumdan yarı-açık duruma geçiş süresi (saniye)
- `failureRatio`: Hata oranı eşiği (0.0-1.0 arası)
- `minimumRequests`: Circuit breaker'ın devreye girmesi için minimum istek sayısı

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

## Özellikler

- Circuit Breaker pattern implementasyonu
- Configurable settings
- Docker desteği
- Health check endpoint'i
- Modüler ve test edilebilir yapı
- Dependency injection
- Fiber web framework entegrasyonu
- Prometheus metrics entegrasyonu
- Rate limiting desteği
- Request logging
- Response caching
- Request validation
- Graceful shutdown
- Swagger/OpenAPI dokümantasyonu

## Middleware'ler

### Logger Middleware
Her HTTP isteği için otomatik logging sağlar:
- HTTP metodu
- İstek yolu
- Yanıt kodu
- İşlem süresi

### Rate Limiter
IP bazlı istek sınırlama:
- Varsayılan limit: 100 istek/saniye
- Burst limit: 10 istek
- Aşım durumunda 429 Too Many Requests yanıtı

### Cache Middleware
GET istekleri için otomatik önbellekleme:
- Varsayılan TTL: 5 dakika
- Path bazlı cache key'leri
- Sadece GET istekleri için aktif

### Validator Middleware
Request body validation:
- Struct tag'leri ile validation kuralları
- Otomatik hata mesajları
- 400 Bad Request yanıtları

## Metrics

Prometheus metrics endpoint'i (`/metrics`):

### Metrik Türleri
- `http_requests_total`: Toplam HTTP istek sayısı
  - Labels: method, path, status
- `circuit_breaker_state`: Circuit breaker durumu
  - Labels: name
  - Values: 0 (Closed), 1 (Half-Open), 2 (Open)

## API Grupları

### API Endpoints (/api)
- `GET /api/`: Ana endpoint
- `GET /api/cached`: Önbellekli endpoint (5 dakika TTL)

### System Endpoints
- `GET /health`: Sağlık kontrolü
- `GET /metrics`: Prometheus metrikleri

## Graceful Shutdown

Uygulama şu sinyalleri yakalayarak graceful shutdown gerçekleştirir:
- SIGTERM
- SIGINT (Ctrl+C)

Shutdown sırasında:
1. Yeni istekleri kabul etmeyi durdurur
2. Mevcut isteklerin tamamlanmasını bekler
3. Kaynakları temizler

## Geliştirme

### Test
```bash
go test ./...
```

### API Dokümantasyonu
Swagger UI: `http://localhost:3000/swagger/index.html`

### Metrics İzleme
1. Prometheus kurulumu:
```bash
docker run -d -p 9090:9090 prom/prometheus
```

2. Grafana ile görselleştirme:
```bash
docker run -d -p 3000:3000 grafana/grafana
```# circuit_breaker_pattern
