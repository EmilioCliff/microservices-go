CREATE TABLE users (
  "id"   BIGSERIAL PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "password" varchar NOT NULL,
  "active" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "update_at" timestamptz NOT NULL DEFAULT now()
);