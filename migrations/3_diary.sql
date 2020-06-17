-- +migrate Up
CREATE TABLE diary (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    year       INTEGER NOT NULL,
    user_id    INTEGER REFERENCES users(id)
);

COMMENT ON TABLE  diary            IS 'Diary';
COMMENT ON COLUMN diary.id         IS 'id';
COMMENT ON COLUMN diary.created_at IS 'create date';
COMMENT ON COLUMN diary.updated_at IS 'update date';
COMMENT ON COLUMN diary.deleted_at IS 'delete date';
COMMENT ON COLUMN diary.year       IS 'year';
COMMENT ON COLUMN diary.user_id    IS 'user';

-- +migrate Down
DROP TABLE diary;