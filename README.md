# Допущения или отход от классики
- В папке docker куча файлов, так как конкретно в этой репе хочу как можно больше технологий заюзать, скоро еще добавится кафка и рэббит для брокеров, так что этих файлов станет еще больше, ниже приведу классическую арху проекта, которую юзаю сам
```
.
├── CHANGELOG.md
├── Dockerfile
├── Makefile
├── README.md
├── .env-docker
├── .env
├── TODO.md
├── api
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
├── cmd
│   └── app
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
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
├── migrations
│   ├── clickhouse
│   │   └── 001_create_dogs_table.sql
│   ├── mysql
│   │   └── 001_create_dogs_table.sql
│   └── postgres
│       └── 001_create_dogs_table.sql
└── pkg
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
