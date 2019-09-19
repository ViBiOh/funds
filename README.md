# Funds

[![Build Status](https://travis-ci.org/ViBiOh/funds.svg?branch=master)](https://travis-ci.org/ViBiOh/funds)
[![codecov](https://codecov.io/gh/ViBiOh/funds/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/funds)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/funds)](https://goreportcard.com/report/github.com/ViBiOh/funds)

## Usage

### API

```bash
Usage of api:
  -address string
        [http] Listen address {API_ADDRESS}
  -cert string
        [http] Certificate file {API_CERT}
  -corsCredentials
        [cors] Access-Control-Allow-Credentials {API_CORS_CREDENTIALS}
  -corsExpose string
        [cors] Access-Control-Expose-Headers {API_CORS_EXPOSE}
  -corsHeaders string
        [cors] Access-Control-Allow-Headers {API_CORS_HEADERS} (default "Content-Type")
  -corsMethods string
        [cors] Access-Control-Allow-Methods {API_CORS_METHODS} (default "GET")
  -corsOrigin string
        [cors] Access-Control-Allow-Origin {API_CORS_ORIGIN} (default "*")
  -csp string
        [owasp] Content-Security-Policy {API_CSP} (default "default-src 'self'; base-uri 'self'")
  -dbHost string
        [db] Host {API_DB_HOST}
  -dbName string
        [db] Name {API_DB_NAME}
  -dbPass string
        [db] Pass {API_DB_PASS}
  -dbPort string
        [db] Port {API_DB_PORT} (default "5432")
  -dbUser string
        [db] User {API_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {API_FRAME_OPTIONS} (default "deny")
  -fundsHour int
        [funds] Hour of running {API_FUNDS_HOUR} (default 8)
  -fundsInterval string
        [funds] Duration between two runs {API_FUNDS_INTERVAL} (default "24h")
  -fundsMaxRetry int
        [funds] Max retry {API_FUNDS_MAX_RETRY} (default 10)
  -fundsMinute int
        [funds] Minute of running {API_FUNDS_MINUTE}
  -fundsOnStart
        [funds] Start scheduler on start {API_FUNDS_ON_START}
  -fundsRetry string
        [funds] Duration between two retries {API_FUNDS_RETRY} (default "10m")
  -fundsTimezone string
        [funds] Timezone of running {API_FUNDS_TIMEZONE} (default "Europe/Paris")
  -hsts
        [owasp] Indicate Strict Transport Security {API_HSTS} (default true)
  -infos string
        [funds] Informations URL {API_INFOS}
  -key string
        [http] Key file {API_KEY}
  -port int
        [http] Listen port {API_PORT} (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics {API_PROMETHEUS_PATH} (default "/metrics")
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) {API_TRACING_AGENT} (default "jaeger:6831")
  -tracingName string
        [tracing] Service name {API_TRACING_NAME}
  -url string
        [alcotest] URL to check {API_URL}
  -userAgent string
        [alcotest] User-Agent for check {API_USER_AGENT} (default "Golang alcotest")
```

### Alert

```bash
Usage of alert:
  -c    Healthcheck (check and exit)
  -dbHost string
        [db] Host {ALERT_DB_HOST}
  -dbName string
        [db] Name {ALERT_DB_NAME}
  -dbPass string
        [db] Pass {ALERT_DB_PASS}
  -dbPort string
        [db] Port {ALERT_DB_PORT} (default "5432")
  -dbUser string
        [db] User {ALERT_DB_USER}
  -infos string
        [funds] Informations URL {ALERT_INFOS}
  -mailerPass string
        [mailer] Mailer Pass {ALERT_MAILER_PASS}
  -mailerURL string
        [mailer] Mailer URL {ALERT_MAILER_URL} (default "https://mailer.vibioh.fr")
  -mailerUser string
        [mailer] Mailer User {ALERT_MAILER_USER}
  -recipients string
        [notifier] Email of notifications recipients {ALERT_RECIPIENTS}
  -schedulerHour int
        [scheduler] Hour of running {ALERT_SCHEDULER_HOUR} (default 8)
  -schedulerInterval string
        [scheduler] Duration between two runs {ALERT_SCHEDULER_INTERVAL} (default "24h")
  -schedulerMaxRetry int
        [scheduler] Max retry {ALERT_SCHEDULER_MAX_RETRY} (default 10)
  -schedulerMinute int
        [scheduler] Minute of running {ALERT_SCHEDULER_MINUTE}
  -schedulerOnStart
        [scheduler] Start scheduler on start {ALERT_SCHEDULER_ON_START}
  -schedulerRetry string
        [scheduler] Duration between two retries {ALERT_SCHEDULER_RETRY} (default "10m")
  -schedulerTimezone string
        [scheduler] Timezone of running {ALERT_SCHEDULER_TIMEZONE} (default "Europe/Paris")
  -score float
        [notifier] Score value to notification when above {ALERT_SCORE} (default 25)
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) {ALERT_TRACING_AGENT} (default "jaeger:6831")
  -tracingName string
        [tracing] Service name {ALERT_TRACING_NAME}
```
