# Funds

[![Build Status](https://travis-ci.org/ViBiOh/funds.svg?branch=master)](https://travis-ci.org/ViBiOh/funds)
[![codecov](https://codecov.io/gh/ViBiOh/funds/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/funds)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/funds)](https://goreportcard.com/report/github.com/ViBiOh/funds)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ViBiOh/funds)](https://dependabot.com)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_funds&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_funds)

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
  -dbSslmode string
        [db] SSL Mode {API_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {API_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {API_FRAME_OPTIONS} (default "deny")
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
  -dbSslmode string
        [db] SSL Mode {ALERT_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {ALERT_DB_USER}
  -infos string
        [funds] Informations URL {ALERT_INFOS}
  -mailerPass string
        [mailer] Pass {ALERT_MAILER_PASS}
  -mailerURL string
        [mailer] URL (an instance of github.com/ViBiOh/mailer) {ALERT_MAILER_URL}
  -mailerUser string
        [mailer] User {ALERT_MAILER_USER}
  -recipients string
        [notifier] Email of notifications recipients {ALERT_RECIPIENTS}
  -score float
        [notifier] Score value to notification when above {ALERT_SCORE} (default 25)
```
