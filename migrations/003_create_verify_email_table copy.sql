CREATE TABLE IF NOT EXISTS verify_email (
  user_id uuid NOT NULL,
  token TEXT UNIQUE NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS verify_email;