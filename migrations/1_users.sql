-- +migrate Up
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    login      VARCHAR(255) UNIQUE NOT NULL,
    password   VARCHAR(255)        NOT NULL
);

COMMENT ON TABLE users IS 'Users';
COMMENT ON COLUMN users.id IS 'id';
COMMENT ON COLUMN users.created_at IS 'create date';
COMMENT ON COLUMN users.updated_at IS 'update date';
COMMENT ON COLUMN users.deleted_at IS 'delete date';
COMMENT ON COLUMN users.login IS 'login';
COMMENT ON COLUMN users.password IS 'password';

-- +migrate Down
DROP TABLE users;