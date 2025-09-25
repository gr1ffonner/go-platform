```
.
├── CHANGELOG.md         # история изменений 
├── Dockerfile           # сборка приложения
├── Makefile             # команды управления
├── README.md            # документация + общая инфа как поднять проект
├── .env-docker          # переменные для докера
├── .env                 # локальные переменные
├── TODO.md              # список задач
├── api                  # API документация
│   ├── docs.go
│   ├── protobuf
│   │   ├── dogs.pb.go
│   │   ├── dogs.proto
│   │   ├── dogs_grpc.pb.go
│   │   ├── health.pb.go
│   │   ├── health.proto
│   │   └── health_grpc.pb.go
│   ├── swagger.json
│   └── swagger.yaml
├── cmd                  # точка входа
│   └── app
│       └── main.go     # приложение
├── docker-compose.yml  # инфраструктура для приложения (бд, брокеры)
├── go.mod              # зависимости
├── go.sum              # хеши зависимостей
├── infra               # конфигурация инфраструктуры мониторинга
│   ├── grafana         # конфиг
│   │   ├── dashboards  # JSON файлы дашбордов
│   │   ├── dashboards.yaml              # автоподключение дашбордов
│   ├── otel-collector  # конфигурация OpenTelemetry Collector
│   │   └── config.yaml # настройки сбора и маршрутизации телеметрии
│   ├── prometheus      # конфигурация Prometheus
│   │   └── scrape_config.yml # настройки сбора метрик
│   └── tempo           # конфигурация Grafana Tempo
│       └── config.yaml # настройки хранения и запросов трейсов
├── internal            # внутренняя логика
│   ├── clients
│   │   ├── rest-client-example
│   │   │   └── client.go
│   │   └── s3
│   │       └── client.go
│   ├── gprc
│   │   ├── dogs-by-breed.go
│   │   ├── healthcheck.go
│   │   ├── middleware.go
│   │   └── server.go
│   ├── handlers
│   │   ├── dog-by-breed.go
│   │   ├── handler.go
│   │   ├── health.go
│   │   ├── middleware.go
│   │   └── router.go
│   ├── models
│   │   └── dogs
│   │       └── models.go
│   ├── services
│   │   └── dogs
│   │       └── service.go
│   └── storages
│       ├── clickhouse
│       │   └── repository.go
│       └── postgresql
│           └── repository.go
├── migrations          # миграции БД
│   ├── clickhouse
│   │   └── 001_create_dogs_table.sql
│   └── postgres
│       └── 001_create_dogs_table.sql
└── pkg                 # переиспользуемые пакеты
    ├── broker
    │   └── nats
    │       └── nats.go
    ├── cache
    │   └── redis
    │       └── redis.go
    ├── config
    │   └── config.go
    ├── db
    │   ├── clickhouse
    │   │   └── clickhouse.go
    │   └── postgre
    │       └── postgres.go
    ├── logger
    │   └── logger.go
    └── utils
        ├── errors.go
        ├── http-utils
        │   ├── errors.go
        │   ├── finalizer.go
        │   └── models.go
        └── utils.go
```

# Комментарии к архитектуре 
- Changelog как правило ведем уже после релиза на прод, для доп трекинга фичей, которые мы релизим, туда же линкуем фичи по возможности [read](https://keepachangelog.com/ru/1.1.0/)