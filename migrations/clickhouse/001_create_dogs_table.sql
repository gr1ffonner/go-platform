-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dogs (
    id String,
    breed String,
    image_url String,
    created_at DateTime64(3) DEFAULT now64(),
    updated_at DateTime64(3) DEFAULT now64()
) ENGINE = MergeTree()
ORDER BY (id, created_at)
PARTITION BY toYYYYMM(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dogs;
-- +goose StatementEnd
