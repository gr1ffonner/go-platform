# Допущения или отход от классики
- В папке docker куча файлов, так как конкретно в этой репе хочу как можно больше технологий заюзать, скоро еще добавится кафка и рэббит для брокеров, так что этих файлов станет еще больше, ниже приведу классическую арху проекта, которую юзаю сам
```
.
├── CHANGELOG.md         # история изменений 
├── Dockerfile           # сборка приложения
├── Makefile             # команды управления
├── README.md            # документация
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
│       └── main.go
├── docker-compose.yml  # инфраструктура для приложения (бд, брокеры)
├── go.mod              # зависимости
├── go.sum              # хеши зависимостей
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
│       ├── mysql
│       │   └── repository.go
│       └── postgresql
│           └── repository.go
├── migrations          # миграции БД
│   ├── clickhouse
│   │   └── 001_create_dogs_table.sql
│   ├── mysql
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
    │   ├── mysql
    │   │   └── mysql.go
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

# Как поднять проект

## Все в докере
```bash
make full-pg # postgresql repo
make full-ch # clickhouse repo
make full-mysql # mysql repo
```

## Инфру в докере, приложение локально
```bash
make up-pg # postgresql repo
make up-ch # clickhouse repo
make up-mysql # mysql repo

make run-pg # postgresql repo
make run-ch # clickhouse repo
make run-mysql # mysql repo
```

# Список технологий 

## Базы данных

### SQL
- Postgresql
- Mysql

### NoSql
- Redis (cache)
- Minio S3 analog (file storage)

## Брокеры сообщений
- Nats
- Kafka (soon)
- RabbitMQ (soon)