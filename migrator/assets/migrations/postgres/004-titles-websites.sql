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
    status       text
);

CREATE TABLE IF NOT EXISTS raw_news_webpages
(
    id         uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at timestamp DEFAULT now(),
    body       text
);

-- +migrate Down
DROP TABLE IF EXISTS titles;
DROP TABLE IF EXISTS raw_news_webpages;