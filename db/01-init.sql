CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS pockets(
	id SERIAL PRIMARY KEY,
	name TEXT,
	category TEXT,
	currency TEXT,
	balance float8
);

CREATE TABLE IF NOT EXISTS transactions(
			id SERIAL PRIMARY KEY,
			source_pid INT,
			dest_pid INT,
			amount float8,
			description TEXT,
			date timestamp,
			status TEXT
);