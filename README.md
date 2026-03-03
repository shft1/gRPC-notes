# gRPC-notes

Сервис заметок на Go с продвинутой реализацией gRPC и protobuf: unary-методы, все виды стриминга, интерцепторы, валидация запросов на уровне `.proto`, ошибки с деталями, HTTP-доступ через gRPC-Gateway, Swagger UI и собственный protoc-плагин.

## 1. Содержание

- [2. О проекте](#2-о-проекте)
- [3. Используемые технологии](#3-используемые-технологии)
- [4. Что реализовано в проекте](#4-что-реализовано-в-проекте)
  - [4.1 Базовый сервис заметок](#41-базовый-сервис-заметок)
  - [4.2 Архитектура и хранение данных](#42-архитектура-и-хранение-данных)
  - [4.3 Логирование и авторизация через интерцепторы](#43-логирование-и-авторизация-через-интерцепторы)
  - [4.4 Настройка параметров gRPC-сервера](#44-настройка-параметров-grpc-сервера)
  - [4.5 Валидация через protobuf](#45-валидация-через-protobuf)
  - [4.6 Ошибки с деталями](#46-ошибки-с-деталями)
  - [4.7 Серверный стриминг](#47-серверный-стриминг)
  - [4.8 Клиентский стриминг](#48-клиентский-стриминг)
  - [4.9 Двунаправленный стриминг](#49-двунаправленный-стриминг)
  - [4.10 Стрим-интерцептор и oneof-ошибки в стриме](#410-стрим-интерцептор-и-oneof-ошибки-в-стриме)
  - [4.11 HTTP, gRPC-Gateway, Swagger и CORS](#411-http-grpc-gateway-swagger-и-cors)
  - [4.12 Два HTTP-пути в проекте и WebSocket для стримов](#412-два-http-пути-в-проекте-и-websocket-для-стримов)
  - [4.13 Кастомный protoc-плагин `protoc-gen-fast-equal`](#413-кастомный-protoc-плагин-protoc-gen-fast-equal)
- [5. API документация](#5-api-документация)
  - [5.1 Unary RPC](#51-unary-rpc)
  - [5.2 Streaming RPC](#52-streaming-rpc)
  - [5.3 REST-маршруты через gRPC-Gateway](#53-rest-маршруты-через-grpc-gateway)
  - [5.4 Форматы сообщений protobuf](#54-форматы-сообщений-protobuf)
  - [5.5 Авторизация, валидация и ошибки](#55-авторизация-валидация-и-ошибки)
- [6. Структура проекта](#6-структура-проекта)
- [7. Локальная развертка](#7-локальная-развертка)
  - [7.1 Требования](#71-требования)
  - [7.2 Установка зависимостей](#72-установка-зависимостей)
  - [7.3 Переменные окружения и за что они отвечают](#73-переменные-окружения-и-за-что-они-отвечают)
  - [7.4 Запуск через Docker Compose](#74-запуск-через-docker-compose)
  - [7.5 Запуск без Docker](#75-запуск-без-docker)
  - [7.6 Генерация кода из `.proto`](#76-генерация-кода-из-proto)
- [8. Команды Makefile](#8-команды-makefile)

## 2. О проекте

**gRPC Notes** — это сервис заметок, где основной API построен на gRPC и protobuf.

Проект показывает не только базовые CRUD-запросы, но и более глубокие возможности:
- цепочки unary и stream интерцепторов;
- server/client/bidi стриминг;
- проверка данных через правила прямо в `.proto`;
- понятная передача ошибок с деталями;
- доступ к API как по gRPC, так и по HTTP через gRPC-Gateway;
- WebSocket-обертка для стримов в web-сценариях;
- собственный плагин для генерации методов сравнения protobuf-структур.

## 3. Используемые технологии

- Go
- gRPC
- Protocol Buffers
- Unary, Server Streaming, Client Streaming, Bidirectional Streaming
- gRPC Interceptors
- Protovalidate (`buf.validate`)
- gRPC-Gateway
- OpenAPI/Swagger
- Chi Router
- CORS middleware
- `grpc-websocket-proxy`
- Docker / Docker Compose
- Makefile
- easyp (generate/lint/breaking)
- custom plugin: `protoc-gen-fast-equal`

## 4. Что реализовано в проекте

### 4.1 Базовый сервис заметок

Реализован protobuf-контракт с одним сервисом `NoteAPI` и методами:
- создание заметки;
- получение одной заметки;
- получение списка заметок;
- удаление заметки.

Для удобной работы с проектом добавлен `Makefile` с командами запуска и генерации.

### 4.2 Архитектура и хранение данных

Бизнес-логика не смешана с хендлерами. Слои разделены:
- `handler` — принимает запросы;
- `usecase` — выполняет бизнес-логику;
- `repository` — работает с данными.

Хранилище in-memory, но потокобезопасное:
- данные лежат в `map`;
- доступ защищен `sync.RWMutex`;
- параллельные чтение/запись работают корректно.

### 4.3 Логирование и авторизация через интерцепторы

Сервер использует `grpc.ChainUnaryInterceptor`.

Реализовано:
- интерцептор логирования unary-запросов:
  - начало вызова;
  - имя метода;
  - длительность выполнения;
  - итог (успех/ошибка);
- интерцептор авторизации unary-запросов:
  - чтение metadata из контекста;
  - проверка `authorization`;
  - ошибка `Unauthenticated`, если токен не прошел проверку.

### 4.4 Настройка параметров gRPC-сервера

Сервер запускается с параметрами соединения (через `grpc.ServerOption`), чтобы лучше контролировать работу под нагрузкой:
- параметры keepalive (`Time`, `Timeout` и связанные настройки);
- ограничение/контроль поведения соединений через конфиг.

### 4.5 Валидация через protobuf

Правила проверки данных описаны прямо в `.proto`:
- обязательные поля;
- минимальная/максимальная длина;
- проверка UUID;
- дополнительные message-правила.

В Go-коде выполняется `protovalidate.Validate(...)`, поэтому не нужно писать много ручных `if`-проверок.

### 4.6 Ошибки с деталями

Ошибки возвращаются не только кодом, но и дополнительными деталями:
- используется gRPC статус с `WithDetails`;
- добавлено protobuf-сообщение `ErrorDetails`;
- клиент может получить не просто «ошибка», а понятную причину.

### 4.7 Серверный стриминг

`SubscribeToEvents`:
- после подключения сразу отправляется health/welcome сообщение;
- при создании новых заметок клиент получает событие;
- при отключении клиента стрим завершается корректно.

### 4.8 Клиентский стриминг

`UploadMetrics`:
- клиент отправляет поток чисел;
- сервер накапливает значения;
- после завершения потока возвращает итоговую сумму.

### 4.9 Двунаправленный стриминг

`Chat`:
- обмен сообщениями в обе стороны в одном стриме;
- чтение и отправка работают асинхронно (горутины);
- используется `correlation_id`, чтобы связывать ответы с исходными сообщениями.

### 4.10 Стрим-интерцептор и oneof-ошибки в стриме

Реализован stream-интерцептор:
- перехватывает сообщения стрима;
- логирует `RecvMsg` и `SendMsg`.

В чате реализована передача бизнес-ошибки без закрытия соединения:
- сообщение использует `oneof`;
- один из вариантов — `google.rpc.Status`.

### 4.11 HTTP, gRPC-Gateway, Swagger и CORS

Сервис доступен по HTTP через gRPC-Gateway:
- в `.proto` добавлены `google.api.http` аннотации;
- сгенерированы gateway-файлы и OpenAPI/Swagger;
- поднят HTTP-сервер для JSON-запросов;
- подключен Swagger UI;
- добавлен CORS middleware для браузера.

### 4.12 Два HTTP-пути в проекте и WebSocket для стримов

В проекте есть два пути работы по HTTP:

1. Самописные HTTP-ручки (`/notes/v1/...`)
- отдельные хендлеры;
- полезно для полного ручного контроля HTTP-логики.

2. gRPC-Gateway (`/gen/v1/...`)
- маршруты автоматически следуют protobuf-контракту;
- добавлена WebSocket-обертка (`wsproxy.WebsocketProxy`) для стриминговых сценариев в web.

### 4.13 Кастомный protoc-плагин `protoc-gen-fast-equal`

В репозитории есть собственный protoc-плагин:
- анализирует protobuf messages;
- генерирует `*_fast_equal.pb.go`;
- добавляет метод `IsEqual(msg *T) bool` для быстрого сравнения.

## 5. API документация

### 5.1 Unary RPC

Сервис: `api.notes.v1.NoteAPI`

- `Create(NoteCreateRequest) returns (Note)`
- `GetByID(NoteIDRequest) returns (Note)`
- `GetMulti(google.protobuf.Empty) returns (NoteList)`
- `DeleteByID(NoteIDRequest) returns (Note)`

### 5.2 Streaming RPC

- `SubscribeToEvents(Empty) returns (stream EventResponse)` — server streaming
- `UploadMetrics(stream MetricRequest) returns (SummaryResponse)` — client streaming
- `Chat(stream Message) returns (stream Message)` — bidirectional streaming

### 5.3 REST-маршруты через gRPC-Gateway

Базовый префикс: `/gen/v1`

- `POST   /gen/v1/notes`
- `GET    /gen/v1/notes/{id}`
- `GET    /gen/v1/notes`
- `DELETE /gen/v1/notes/{id}`
- `GET    /gen/v1/notes/subscribe`
- `POST   /gen/v1/notes/metrics/upload`
- `POST   /gen/v1/notes/chat`

Swagger:
- UI: `/gen/v1/swagger/`
- спецификации: `/gen/v1/swagger/specs/...`

### 5.4 Форматы сообщений protobuf

Основные сообщения API:
- `NoteCreateRequest`
- `NoteIDRequest`
- `Note`
- `NoteList`
- `EventResponse` (`oneof` для события/health)
- `MetricRequest`
- `SummaryResponse`
- `Message` (`correlation_id` + `oneof payload`)
- `ErrorDetails` (детали ошибок)

### 5.5 Авторизация, валидация и ошибки

- Авторизация проверяется по `authorization` metadata/header.
- Валидация запросов выполняется по правилам в `.proto`.
- При ошибках возвращаются корректные gRPC-коды (`InvalidArgument`, `NotFound`, `FailedPrecondition`, `Unauthenticated`, `Internal`).
- В ряде случаев передаются дополнительные details для клиента.

## 6. Структура проекта

- `cmd/app` — запуск gRPC-сервера
- `cmd/client` — запуск HTTP-слоя
- `api/notes/v1` — protobuf-контракты
- `pkg/api/notes/v1` — сгенерированный код protobuf/gRPC/gateway
- `internal/app` — backend-слои (config, handler, usecase, repository, middleware, server)
- `internal/client` — HTTP/gateway-слои (router, handler, gateway, middleware, server)
- `docs/api/notes/v1` — Swagger JSON
- `static/swagger` — встроенные swagger-статические файлы
- `protoc-gen-fast-equal` — исходники кастомного плагина
- `deploy/app`, `deploy/client` — docker-конфигурация

## 7. Локальная развертка

### 7.1 Требования

- Go 1.23+
- Docker + Docker Compose
- protoc (если нужна ручная генерация)

### 7.2 Установка зависимостей

После клонирования репозитория:

```bash
go mod download
go mod tidy
```

Установка инструментов генерации:

```bash
make bin-deps
```

### 7.3 Переменные окружения и за что они отвечают

Используются два env-файла:
- `deploy/app/.env` — настройки gRPC backend
- `deploy/client/.env` — настройки HTTP/gateway

Пример `deploy/app/.env`:

```env
GRPC_PORT=50051
MAX_CONNECTION_IDLE=15s
MAX_CONNECTION_AGE=30s
MAX_CONNECTION_AGE_GRACE=5s
TIME=10s
TIMEOUT=5s
CAPACITY=10
```

Что означает:
- `GRPC_PORT` — порт gRPC сервера
- `MAX_CONNECTION_IDLE` — время простоя соединения
- `MAX_CONNECTION_AGE` — максимальное время жизни соединения
- `MAX_CONNECTION_AGE_GRACE` — мягкий период завершения соединения
- `TIME` — интервал keepalive ping
- `TIMEOUT` — ожидание ответа на ping
- `CAPACITY` — емкость буфера event bus

Пример `deploy/client/.env`:

```env
CLIENT_HOST=0.0.0.0
CLIENT_PORT=8080
GRPC_HOST=service-notes
GRPC_PORT=50051
```

Что означает:
- `CLIENT_HOST` — адрес HTTP сервера
- `CLIENT_PORT` — порт HTTP сервера
- `GRPC_HOST` — адрес gRPC backend для подключения
- `GRPC_PORT` — порт gRPC backend

### 7.4 Запуск через Docker Compose

1. Создайте/заполните:
- `deploy/app/.env`
- `deploy/client/.env`

2. Запустите сервисы:

```bash
make run
```

3. Доступность:
- gRPC: `localhost:50051`
- HTTP: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/gen/v1/swagger/`

### 7.5 Запуск без Docker

1. Подготовьте переменные окружения (через `.env` или export).
2. Запустите gRPC backend:

```bash
go run ./cmd/app
```

3. Запустите HTTP слой:

```bash
go run ./cmd/client
```

### 7.6 Генерация кода из `.proto`

```bash
make bin-deps
make generate
```

Дополнительно:

```bash
make lint
make breaking
```

## 8. Команды Makefile

- `make bin-deps` — установить инструменты генерации
- `make generate` — сгенерировать protobuf/gRPC/gateway/swagger код
- `make lint` — проверить proto-контракты
- `make breaking` — проверить совместимость контрактов
- `make run` — запустить сервисы в Docker
- `make down` — остановить сервисы
