# Funds

[![Build Status](https://travis-ci.org/ViBiOh/funds.svg?branch=master)](https://travis-ci.org/ViBiOh/funds)
[![codecov](https://codecov.io/gh/ViBiOh/funds/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/funds)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/funds)](https://goreportcard.com/report/github.com/ViBiOh/funds)

## Usage

### API

```bash
Usage of api:
  -cert string
        [http] Certificate file
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
        [db] Host
  -dbName string
        [db] Name
  -dbPass string
        [db] Pass
  -dbPort string
        [db] Port (default "5432")
  -dbUser string
        [db] User
  -frameOptions string
        [owasp] X-Frame-Options (default "deny")
  -fundsHour int
        [funds] Hour of running (default 8)
  -fundsInterval string
        [funds] Duration between two runs (default "24h")
  -fundsMaxRetry int
        [funds] Max retry (default 10)
  -fundsMinute int
        [funds] Minute of running
  -fundsOnStart
        [funds] Start scheduler on start
  -fundsRetry string
        [funds] Duration between two retries (default "10m")
  -fundsTimezone string
        [funds] Timezone of running (default "Europe/Paris")
  -hsts
        [owasp] Indicate Strict Transport Security (default true)
  -infos string
        [funds] Informations URL
  -key string
        [http] Key file
  -port int
        [http] Listen port (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics (default "/metrics")
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
        [tracing] Service name
  -url string
        [alcotest] URL to check
  -userAgent string
        [alcotest] User-Agent for check (default "Golang alcotest")
```

### Alert

```bash
Usage of alert:
  -c    Healthcheck (check and exit)
  -dbHost string
        [db] Host
  -dbName string
        [db] Name
  -dbPass string
        [db] Pass
  -dbPort string
        [db] Port (default "5432")
  -dbUser string
        [db] User
  -infos string
        [funds] Informations URL
  -mailerPass string
        [mailer] Mailer Pass
  -mailerURL string
        [mailer] Mailer URL (default "https://mailer.vibioh.fr")
  -mailerUser string
        [mailer] Mailer User
  -recipients string
        Email of notifications recipients
  -schedulerHour int
        [scheduler] Hour of running (default 8)
  -schedulerInterval string
        [scheduler] Duration between two runs (default "24h")
  -schedulerMaxRetry int
        [scheduler] Max retry (default 10)
  -schedulerMinute int
        [scheduler] Minute of running
  -schedulerOnStart
        [scheduler] Start scheduler on start
  -schedulerRetry string
        [scheduler] Duration between two retries (default "10m")
  -schedulerTimezone string
        [scheduler] Timezone of running (default "Europe/Paris")
  -score float
        Score value to notification when above (default 25)
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
        [tracing] Service name
```
