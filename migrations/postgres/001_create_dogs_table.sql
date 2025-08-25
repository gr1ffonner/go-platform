-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dogs (
    id SERIAL PRIMARY KEY,
    breed VARCHAR(100) NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dogs;
-- +goose StatementEnd
