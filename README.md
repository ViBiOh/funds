# Funds

[![Build Status](https://travis-ci.org/ViBiOh/funds.svg?branch=master)](https://travis-ci.org/ViBiOh/funds)
[![codecov](https://codecov.io/gh/ViBiOh/funds/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/funds)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/funds)](https://goreportcard.com/report/github.com/ViBiOh/funds)

## Postgres installation

```bash
export FUNDS_DATABASE_DIR=`realpath ./data`
export FUNDS_DATABASE_PASS=password

mkdir ${FUNDS_DATABASE_DIR}
sudo chown -R 70:70 ${FUNDS_DATABASE_DIR}

docker-compose -p funds -f docker-compose.db.yml up -d
```

## Usage

### API

```bash
Usage of api:
  -corsCredentials
      [cors] Access-Control-Allow-Credentials
  -corsExpose string
      [cors] Access-Control-Expose-Headers
  -corsHeaders string
      [cors] Access-Control-Allow-Headers (default "Content-Type")
  -corsMethods string
      [cors] Access-Control-Allow-Methods (default "GET")
  -corsOrigin string
      [cors] Access-Control-Allow-Origin (default "*")
  -csp string
      [owasp] Content-Security-Policy (default "default-src 'self'; base-uri 'self'")
  -dbHost string
      [database] Host
  -dbName string
      [database] Name
  -dbPass string
      [database] Pass
  -dbPort string
      [database] Port (default "5432")
  -dbUser string
      [database] User
  -frameOptions string
      [owasp] X-Frame-Options (default "deny")
  -hsts
      [owasp] Indicate Strict Transport Security (default true)
  -infos string
      [funds] Informations URL
  -port int
      Listen port (default 1080)
  -prometheusPath string
      [prometheus] Path for exposing metrics (default "/metrics")
  -tls
      Serve TLS content (default true)
  -tlsCert string
      [tls] PEM Certificate file
  -tlsHosts string
      [tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
      [tls] PEM Key file
  -tlsOrganization string
      [tls] Self-signed certificate organization (default "ViBiOh")
  -tracingAgent string
      [opentracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
      [opentracing] Service name
  -url string
      [health] URL to check
  -userAgent string
      [health] User-Agent for check (default "Golang alcotest")
```

### Alert

```bash
Usage of alert:
  -c  Healthcheck (check and exit)
  -dbHost string
      [database] Host
  -dbName string
      [database] Name
  -dbPass string
      [database] Pass
  -dbPort string
      [database] Port (default "5432")
  -dbUser string
      [database] User
  -hour int
      Hour of day for sending notifications (default 8)
  -infos string
      [funds] Informations URL
  -mailerPass string
      Mailer Pass
  -mailerURL string
      Mailer URL
  -mailerUser string
      Mailer User
  -minute int
      Minute of hour for sending notifications
  -recipients string
      Email of notifications recipients
  -score float
      Score value to notification when above (default 25)
  -timezone string
      Timezone (default "Europe/Paris")
  -tracingAgent string
      [opentracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
      [opentracing] Service name
```
