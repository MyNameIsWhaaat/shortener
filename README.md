# URL Shortener Service

Полнофункциональный микросервис для сокращения длинных URL-адресов с встроенной аналитикой, кэшированием популярных ссылок и веб-интерфейсом.

## Возможности

- Создание сокращённых ссылок с автоматической генерацией кодов или кастомными именами
- Редирект с отслеживанием переходов (User-Agent, IP, Referer, timestamp)
- Аналитика по дням, месяцам и типам устройств
- Redis кэширование для популярных ссылок с использованием Sorted Sets
- Веб-интерфейс для управления ссылками и просмотра аналитики
- RESTful API для интеграции в другие приложения

## Требования

- Go 1.24.2 или выше
- Docker и Docker Compose
- Или локально: PostgreSQL 16+, Redis 7+

## Быстрый старт

### С Docker Compose

```bash
cd shortener
docker-compose up --build
```

Приложение доступно на http://localhost:8080

### Локально

1. Установка зависимостей:
```bash
go mod download
```

2. Запуск PostgreSQL и Redis:
```bash
docker-compose up postgres redis migrate
```

3. Запуск приложения:
```bash
go run ./cmd/shortener/main.go
```

## API Endpoints

### POST /api/shorten

Создание сокращённой ссылки.

Запрос:
```json
{
  "url": "https://example.com/very/long/url",
  "custom_alias": "my-link"
}
```

Ответ:
```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:8080/s/abc123",
  "original_url": "https://example.com/very/long/url"
}
```

### GET /s/{short_code}

Редирект на оригинальный URL с сохранением информации о переходе.

### GET /api/analytics/{short_code}

Получение полной аналитики по ссылке.

Ответ:
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com/very/long/url",
  "created_at": "2026-02-20T10:30:00Z",
  "total_clicks": 42,
  "daily_stats": {
    "2026-02-20": 15,
    "2026-02-19": 27
  },
  "devices": {
    "Desktop": 25,
    "Mobile": 17
  },
  "recent_clicks": [
    {
      "user_agent": "Mozilla/5.0...",
      "ip": "192.168.1.1",
      "referer": "https://google.com",
      "created_at": "2026-02-20T10:45:00Z"
    }
  ]
}
```

### GET /api/urls

Получение всех ссылок с пагинацией.

Параметры:
- limit: количество результатов (по умолчанию 20)

### GET /api/urls/popular

Получение популярных ссылок из Redis кэша.

Параметры:
- limit: количество результатов (по умолчанию 10)

## Веб-интерфейс

Доступен на http://localhost:8080/

Функциональность:
- Форма для создания новых сокращённых ссылок
- Таблица всех ссылок с быстрым доступом
- Модальное окно аналитики с подробной статистикой
- Визуализация данных по дням и устройствам
- Список последних переходов

## Архитектура

Структура проекта:
```
internal/
  api/              - HTTP маршруты
  httpapi/          - HTTP handlers и middleware
  service/          - Бизнес-логика
  store/            - Работа с БД
  cache/            - Redis интеграция
  config/           - Конфигурация
  domain/           - Доменные модели
  ui/               - Веб-интерфейс

cmd/shortener/     - Точка входа
migrations/        - SQL миграции
```

Слои приложения:
1. Domain - доменные модели (URL, ClickEvent, AnalyticsResponse)
2. Store - операции с PostgreSQL (CRUD)
3. Service - бизнес-логика и валидация
4. Cache - Redis интеграция с fallback
5. API - HTTP handlers и маршруты

## Переменные окружения

```
PORT=8080
BASE_URL=http://localhost:8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=shortener
DB_PASSWORD=password
DB_NAME=shortener
DB_SSL_MODE=disable

REDIS_ADDR=redis:6379
CACHE_TTL=24h

SHORT_CODE_LENGTH=6
```

## Тестирование

Запуск всех тестов:
```bash
go test ./...
```

Запуск с покрытием:
```bash
go test -v -cover ./...
```

Статистика тестов:
- 40+ unit тестов
- Полное покрытие основных пакетов (config, domain, service, store, httpapi)

## Примеры использования

Создание ссылки:
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com","custom_alias":"gh"}'
```

Получение аналитики:
```bash
curl http://localhost:8080/api/analytics/gh
```

Проверка Redis кэша:
```bash
docker-compose exec redis redis-cli KEYS "*"
```

Health check:
```bash
curl http://localhost:8080/health
```

## Структура БД

Таблица urls:
```sql
id bigint PRIMARY KEY
short_code varchar(50) UNIQUE NOT NULL
original_url text NOT NULL
custom_alias varchar(50)
created_at timestamptz NOT NULL DEFAULT NOW()
clicks bigint NOT NULL DEFAULT 0
```

Таблица click_events:
```sql
id bigint PRIMARY KEY
short_code varchar(50)
user_agent text
ip varchar(45)
referer text
created_at timestamptz NOT NULL DEFAULT NOW()
```

## Кэширование

Redis используется для:
- Кэширования популярных ссылок (Sorted Set с рейтингом)
- Быстрого доступа к часто используемым URL
- Инкрементирования счётчиков популярности

При недоступности Redis приложение продолжает работать с использованием NoOpCache.

## Развертывание

Production развертывание:
```bash
docker-compose up -d
```

Сервис является stateless и может масштабироваться горизонтально с общей БД и Redis.

## Зависимости

- github.com/gorilla/mux - HTTP маршруты
- github.com/redis/go-redis/v9 - Redis клиент
- github.com/lib/pq - PostgreSQL драйвер

## Дополнительно реализованное

- Unit тесты (40+ тестов)
- Graceful shutdown
- Clean Architecture с разделением слоёв
- Логирование операций
- Валидация URL и кастомных кодов
- Fallback при отключении Redis
- Миграции БД

## Версия и технологический стек

- Go: 1.24.2
- PostgreSQL: 16
- Redis: 7-alpine
- База данных миграций: migrate CLI
