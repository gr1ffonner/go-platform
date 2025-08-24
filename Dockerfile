FROM golang:1.22 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

COPY ./ ./

# Install swag and build the application
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
	go mod download && \
    swag fmt && swag init --pdl=1 -g cmd/app/main.go -o api/ && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app && \

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/api ./api

EXPOSE 8080

CMD ["./main"]
