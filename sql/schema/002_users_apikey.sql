-- +goose Up
ALTER TABLE users
ADD COLUMN apikey VARCHAR(64) UNIQUE NOT NULL 
DEFAULT encode(sha256(random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users
DROP COLUMN apikey;