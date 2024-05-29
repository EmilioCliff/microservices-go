CREATE TABLE logs (
    "id" bigint PRIMARY KEY,
    "email" varchar NOT NULL,
    "data" text NOT NULL,
    "user_ip" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "logged_at" timestamptz NOT NULL DEFAULT now()
)