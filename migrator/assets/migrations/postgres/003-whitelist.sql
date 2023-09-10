-- +migrate Up
ALTER TABLE users
    ADD COLUMN role text;
ALTER TABLE users
    ADD CONSTRAINT username_unq UNIQUE (username);

CREATE TABLE IF NOT EXISTS whitelist
(
    id       uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    username text UNIQUE REFERENCES users (username),
    token    uuid
);

-- +migrate Down
ALTER TABLE users
    DROP COLUMN role;

DROP TABLE IF EXISTS whitelist;

ALTER TABLE users
    DROP CONSTRAINT username_unq;