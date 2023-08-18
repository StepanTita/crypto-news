-- +migrate Up
CREATE TABLE IF NOT EXISTS titles
(
    id           uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at   timestamp DEFAULT now(),
    updated_at   timestamp,
    title        text,
    summary      text,
    hash         text UNIQUE,
    url          text,
    release_date date,
    status       text,
    source       text
);

CREATE TABLE IF NOT EXISTS raw_news
(
    id         uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at timestamp DEFAULT now(),
    title_id   uuid REFERENCES titles (id) NOT NULL,
    body       text
);

-- +migrate Down
DROP TABLE IF EXISTS titles;
DROP TABLE IF EXISTS raw_news;