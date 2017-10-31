-- Cleaning
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS funds;

DROP INDEX IF EXISTS alerts_id;
DROP INDEX IF EXISTS alerts_isin;
DROP INDEX IF EXISTS funds_isin;

DROP SEQUENCE IF EXISTS alerts_id_seq;
DROP TYPE IF EXISTS alert_type;

-- Funds
CREATE TABLE funds (
  isin TEXT NOT NULL,
  label TEXT NOT NULL,
  score NUMERIC(5,2) NOT NULL,
  creation_date TIMESTAMP DEFAULT now(),
  update_date TIMESTAMP
);

CREATE UNIQUE INDEX funds_isin ON funds (isin);

-- Alerts
CREATE TYPE alert_type AS ENUM ('above', 'below');

CREATE SEQUENCE alerts_id_seq;

CREATE TABLE alerts (
  id INTEGER DEFAULT nextval('alerts_id_seq') NOT NULL,
  isin TEXT NOT NULL REFERENCES funds(isin),
  score NUMERIC(5,2) NOT NULL,
  type alert_type NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX alerts_id ON alerts (id);
CREATE INDEX alerts_isin ON alerts (isin);
