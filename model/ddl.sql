DROP TABLE IF EXISTS alerts;
DROP INDEX IF EXISTS funds_isin;
DROP TABLE IF EXISTS funds;
DROP TYPE IF EXISTS alert_type;

CREATE TABLE funds (
  isin TEXT NOT NULL,
  score NUMERIC(5,2) NOT NULL,
  creation_date TIMESTAMP DEFAULT now(),
  update_date TIMESTAMP
);

CREATE UNIQUE INDEX funds_isin ON funds (isin);

CREATE TYPE alert_type AS ENUM ('above', 'below');

CREATE TABLE alerts (
  isin TEXT NOT NULL REFERENCES funds(isin),
  type alert_type NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);