CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  average_rating NUMERIC(2,1) DEFAULT 0.0 NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  email TEXT UNIQUE NOT NULL,
  email_verified_at TIMESTAMPTZ,
  favourites TEXT[], 
  firstname TEXT  NOT NULL,
  image TEXT DEFAULT '' NOT NULL,
  lastname TEXT NOT NULL,
  password TEXT NOT NULL,
  phone_verified_at TIMESTAMPTZ,
  otp_secret_key TEXT DEFAULT '' NOT NULL,
  trips_count INT DEFAULT 0 NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;