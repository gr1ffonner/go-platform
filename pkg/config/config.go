package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	Redis           RedisConfig
	NATS            NATSConfig
	Logger          Logger
	S3              S3
	DogAPI          DogAPIConfig
	MetricsProvider MetricsProviderConfig
}

type ServerConfig struct {
	HTTPPort string `env:"SERVER_PORT" env-default:"8080"`
	GRPCPort string `env:"GRPC_PORT" env-default:"50051"`
	Storage  string `env:"STORAGE" env-default:"postgres"`
}

type S3 struct {
	KeyID              string `env:"S3_STORAGE_KEY"`             // access key ID
	KeySecret          string `env:"S3_STORAGE_SECRET"`          // secret key
	Bucket             string `env:"S3_STORAGE_BUCKET"`          // bucket name
	BaseEndpoint       string `env:"S3_STORAGE_ENDPOINT"`        // S3 endpoint (например, http://localhost:9000)
	BasePublicEndpoint string `env:"S3_STORAGE_ENDPOINT_PUBLIC"` // S3 public endpoint
	Region             string `env:"S3_STORAGE_REGION"`          // region (например, us-east-1)
}

type DogAPIConfig struct {
	BaseURL string `env:"DOG_API_BASE_URL" env-default:"www.example.com"`
}

type Logger struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

type DatabaseConfig struct {
	PostgresDSN   string `env:"POSTGRES_DSN"`
	MySQLDSN      string `env:"MYSQL_DSN"`
	ClickHouseDSN string `env:"CLICKHOUSE_DSN"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

type NATSConfig struct {
	URL string `env:"NATS_URL" env-default:"nats://localhost:4222"`
}

type MetricsProviderConfig struct {
	ServiceName    string `env:"OTEL_SERVICE_NAME" env-default:"go-platform"`
	ServiceVersion string `env:"OTEL_SERVICE_VERSION" env-default:"1.0.0"`
	Environment    string `env:"ENVIRONMENT" env-default:"development"`
	OTLPEndpoint   string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" env-default:"localhost:4317"`
	Insecure       bool   `env:"OTEL_EXPORTER_OTLP_INSECURE" env-default:"true"`
	PrometheusPort string `env:"PROMETHEUS_PORT" env-default:"9090"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
