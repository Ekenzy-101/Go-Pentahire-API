CREATE TABLE IF NOT EXISTS vehicles (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  address TEXT NOT NULL,
  average_rating NUMERIC(2,1) DEFAULT 0.0 NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  is_rented BOOLEAN DEFAULT false NOT NULL,
  image TEXT DEFAULT '' NOT NULL,
  make TEXT NOT NULL,
  name TEXT NOT NULL,
  location POINT NOT NULL,
  rental_fee INT  NOT NULL,
  reviews_count INT DEFAULT 0 NOT NULL,
  user_id uuid NOT NULL,
  trips_count INT DEFAULT 0 NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS vehicles;
