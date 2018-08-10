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
