-- +migrate Up
ALTER TABLE news ADD COLUMN locale text;

-- +migrate Down
ALTER TABLE news DROP COLUMN locale;