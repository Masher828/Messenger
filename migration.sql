SELECT 'CREATE DATABASE social_db' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'social_db')\gexec

\c social_db;

DROP SEQUENCE IF EXISTS hibernate_sequence;
DROP INDEX IF EXISTS idx_social_user_email;
DROP TABLE IF EXISTS social_user;

CREATE SEQUENCE hibernate_sequence
INCREMENT BY 1
MINVALUE 1
MAXVALUE 9223372036854775807
START 1;

CREATE TABLE social_user (
id           int8          NOT NULL,
first_name    varchar(255) NOT NULL,
last_name     varchar(255) DEFAULT '',
email         varchar(255) UNIQUE NOT NULL,
password      varchar(255) NOT NULL,
contact       varchar(15),
country_code  varchar(10),
gender        varchar(255),
country       varchar(255),
date_of_birth timestamptz,
date_created  timestamptz NOT NULL,
last_updated  timestamptz NOT NULL,
last_login    timestamptz NOT NULL,
CONSTRAINT social_user_pkey_id    PRIMARY KEY (id)
);

CREATE INDEX idx_social_user_email on social_user USING btree(email);
