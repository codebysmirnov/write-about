-- +migrate Up
CREATE TABLE person
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255),
    phone      VARCHAR(32),
    email      VARCHAR(56),
    user_id    INTEGER REFERENCES users (id)
);

COMMENT ON TABLE person IS 'Person';
COMMENT ON COLUMN person.created_at IS 'created date';
COMMENT ON COLUMN person.updated_at IS 'update date';
COMMENT ON COLUMN person.deleted_at IS 'deleted date';
COMMENT ON COLUMN person.first_name IS 'first name';
COMMENT ON COLUMN person.last_name IS 'last name';
COMMENT ON COLUMN person.phone IS 'phone';
COMMENT ON COLUMN person.email IS 'email';
COMMENT ON COLUMN person.user_id IS 'user';


-- +migrate Down
DROP TABLE person;