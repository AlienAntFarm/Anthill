CREATE SCHEMA IF NOT EXISTS anthive;

-- Set default search_path to schema
SET search_path TO anthive,public;

DROP TYPE IF EXISTS state;
CREATE TYPE state AS enum ('NEW', 'PENDING', 'FINISH', 'ERROR');

-- Creation of tables
CREATE TABLE IF NOT EXISTS antling (
  id serial primary key
);

CREATE TABLE IF NOT EXISTS image (
  id serial primary key,
  archive varchar(10),
  command text[],
  environment text[],
  cwd text DEFAULT '/',
  hostname text
);

CREATE TABLE IF NOT EXISTS job (
  id serial primary key,
  state state DEFAULT 'NEW',
  fk_antling integer references antling (id),
  fk_image integer references image (id),
  command text[],
  environment text[],
  cwd text
);
