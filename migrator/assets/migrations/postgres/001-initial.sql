-- +migrate Up

CREATE TABLE IF NOT EXISTS users
(
    id         uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    first_name text,
    last_name  text,
    username   text,
    platform   text
);

CREATE TABLE IF NOT EXISTS news
(
    id              uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at      timestamp DEFAULT now(),
    updated_at      timestamp,
    published_at    timestamp,
    media           jsonb,
    url             text,
    source          text,
    original_source text
);

CREATE TABLE IF NOT EXISTS coins
(
    code  text PRIMARY KEY,
    title text NOT NULL,
    slug  text NOT NULL
);

CREATE TABLE IF NOT EXISTS news_coins
(
    id      uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    code    text NOT NULL REFERENCES coins (code),
    news_id uuid NOT NULL REFERENCES news (id),
    UNIQUE (code, news_id)
);

CREATE TABLE IF NOT EXISTS channels
(
    channel_id bigint PRIMARY KEY,
    created_at timestamp DEFAULT now(),
    platform   text,
    priority   bigserial NOT NULL
);

CREATE TABLE IF NOT EXISTS news_channels
(
    id         uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    channel_id bigint REFERENCES channels (channel_id),
    news_id    uuid NOT NULL REFERENCES news (id),
    UNIQUE (channel_id, news_id)
);

CREATE TABLE IF NOT EXISTS preferences_channel_coins
(
    channel_id bigint REFERENCES channels (channel_id),
    coin_code  text NOT NULL REFERENCES coins (code),
    PRIMARY KEY (channel_id, coin_code)
);

-- +migrate Down
DROP TABLE IF EXISTS preferences_regions;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS news_channels;
DROP TABLE IF EXISTS news_coins;
DROP TABLE IF EXISTS channels;
DROP TABLE IF EXISTS coins;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS news;