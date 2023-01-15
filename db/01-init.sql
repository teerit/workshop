CREATE SEQUENCE IF NOT EXISTS account_id;
CREATE SEQUENCE IF NOT EXISTS pocket_id;
CREATE SEQUENCE IF NOT EXISTS transaction_id;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE "pockets" (
    "id" int4 NOT NULL DEFAULT nextval('pocket_id'::regclass),
    "name" TEXT NOT NULL DEFAULT 0,
    "category" TEXT NOT NULL DEFAULT 0,
    "currency" TEXT NOT NULL DEFAULT 0,
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE "transactions" (
    "id" int4 NOT NULL DEFAULT nextval('transaction_id'::regclass),
    "source_pid" int4 NOT NULL DEFAULT 0,
    "dest_pid" int4 NOT NULL DEFAULT 0,
    "amount" float8 NOT NULL DEFAULT 0,
    "description" TEXT NOT NULL DEFAULT 0,
    "date" timestamp NOT NULL,
    "status" TEXT NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);
