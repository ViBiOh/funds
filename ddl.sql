-- Cleaning
DROP INDEX IF EXISTS funds.alerts_id;
DROP INDEX IF EXISTS funds.alerts_isin;
DROP INDEX IF EXISTS funds.funds_isin;

DROP TABLE IF EXISTS funds.alerts;
DROP TABLE IF EXISTS funds.funds;

DROP SEQUENCE IF EXISTS funds.alerts_id_seq;
DROP TYPE IF EXISTS funds.alert_type;

DROP SCHEMA IF EXISTS funds

-- schema
CREATE SCHEMA funds;

-- Funds
CREATE TABLE funds.funds (
  isin TEXT NOT NULL,
  label TEXT NOT NULL,
  score NUMERIC(5,2) NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE DEFAULT now(),
  update_date TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX funds_isin ON funds.funds(isin);

-- Alerts
CREATE TYPE funds.alert_type AS ENUM ('above', 'below');

CREATE SEQUENCE funds.alerts_id_seq;
CREATE TABLE funds.alerts (
  id INTEGER DEFAULT nextval('funds.alerts_id_seq') NOT NULL,
  isin TEXT NOT NULL REFERENCES funds.funds(isin),
  score NUMERIC(5,2) NOT NULL,
  type alert_type NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNIQUE INDEX alerts_id ON funds.alerts(id);
CREATE INDEX alerts_isin ON funds.alerts(isin);
