# Web Monitor

## Description

Web Monitor periodically monitors web URLs and records their metrics.

## Prerequisites

- Go
- Docker

## Configuration

```bash
$ mkdir -p data/postgres
```

## Run

```bash
$ docker compose up
```

### Create Web Url for monitoring
Call POST [http://localhost:8080/web-url](http://localhost:8080/web-url)

Body:

```json
{
    "url": "http://example.com",
    "interval": 2,
    "regexPattern": "xyz"
}