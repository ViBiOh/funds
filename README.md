# Funds

[![Build](https://github.com/ViBiOh/funds/workflows/Build/badge.svg)](https://github.com/ViBiOh/funds/actions)
[![codecov](https://codecov.io/gh/ViBiOh/funds/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/funds)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/funds)](https://goreportcard.com/report/github.com/ViBiOh/funds)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_funds&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_funds)

## Usage

### API

```bash
Usage of api:
  -address string
        [server] Listen address {API_ADDRESS}
  -cert string
        [server] Certificate file {API_CERT}
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
  -dbPort uint
        [db] Port {API_DB_PORT} (default 5432)
  -dbSslmode string
        [db] SSL Mode {API_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {API_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {API_FRAME_OPTIONS} (default "deny")
  -graceDuration string
        [http] Grace duration when SIGTERM received {API_GRACE_DURATION} (default "30s")
  -hsts
        [owasp] Indicate Strict Transport Security {API_HSTS} (default true)
  -idleTimeout string
        [server] Idle Timeout {API_IDLE_TIMEOUT} (default "2m")
  -infos string
        [funds] Informations URL {API_INFOS}
  -key string
        [server] Key file {API_KEY}
  -loggerJson
        [logger] Log format as JSON {API_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {API_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {API_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {API_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {API_LOGGER_TIME_KEY} (default "time")
  -okStatus int
        [http] Healthy HTTP Status code {API_OK_STATUS} (default 204)
  -port uint
        [server] Listen port {API_PORT} (default 1080)
  -prometheusAddress string
        [prometheus] Listen address {API_PROMETHEUS_ADDRESS}
  -prometheusCert string
        [prometheus] Certificate file {API_PROMETHEUS_CERT}
  -prometheusIdleTimeout string
        [prometheus] Idle Timeout {API_PROMETHEUS_IDLE_TIMEOUT} (default "10s")
  -prometheusIgnore string
        [prometheus] Ignored path prefixes for metrics, comma separated {API_PROMETHEUS_IGNORE} (default "/ready")
  -prometheusKey string
        [prometheus] Key file {API_PROMETHEUS_KEY}
  -prometheusPort uint
        [prometheus] Listen port {API_PROMETHEUS_PORT} (default 9090)
  -prometheusReadTimeout string
        [prometheus] Read Timeout {API_PROMETHEUS_READ_TIMEOUT} (default "5s")
  -prometheusShutdownTimeout string
        [prometheus] Shutdown Timeout {API_PROMETHEUS_SHUTDOWN_TIMEOUT} (default "5s")
  -prometheusWriteTimeout string
        [prometheus] Write Timeout {API_PROMETHEUS_WRITE_TIMEOUT} (default "10s")
  -readTimeout string
        [server] Read Timeout {API_READ_TIMEOUT} (default "5s")
  -shutdownTimeout string
        [server] Shutdown Timeout {API_SHUTDOWN_TIMEOUT} (default "10s")
  -url string
        [alcotest] URL to check {API_URL}
  -userAgent string
        [alcotest] User-Agent for check {API_USER_AGENT} (default "Alcotest")
  -writeTimeout string
        [server] Write Timeout {API_WRITE_TIMEOUT} (default "10s")
```

### Notifier

```bash
Usage of notifier:
  -c    Healthcheck (check and exit)
  -cron
        [notifier] Start as a cron {NOTIFIER_CRON}
  -dbHost string
        [db] Host {NOTIFIER_DB_HOST}
  -dbName string
        [db] Name {NOTIFIER_DB_NAME}
  -dbPass string
        [db] Pass {NOTIFIER_DB_PASS}
  -dbPort uint
        [db] Port {NOTIFIER_DB_PORT} (default 5432)
  -dbSslmode string
        [db] SSL Mode {NOTIFIER_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {NOTIFIER_DB_USER}
  -infos string
        [funds] Informations URL {NOTIFIER_INFOS}
  -mailerName string
        [mailer] HTTP Username or AMQP Exchange name {NOTIFIER_MAILER_NAME} (default "mailer")
  -mailerPassword string
        [mailer] HTTP Pass {NOTIFIER_MAILER_PASSWORD}
  -mailerURL string
        [mailer] URL (https?:// or amqps?://) {NOTIFIER_MAILER_URL}
  -recipients string
        [notifier] Email of notifications recipients {NOTIFIER_RECIPIENTS}
  -score float
        [notifier] Score value to notification when above {NOTIFIER_SCORE} (default 25)
```
