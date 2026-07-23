# Orderly

Order Processing API — Go + Gin, PostgreSQL, Redis.

## Status

Work in progress.

## Stack

- Go / Gin
- PostgreSQL (golang-migrate)
- Redis (caching, rate limiting)
- Docker Compose
- GitHub Actions

## Getting started

```bash
docker compose up -d
export DATABASE_URL="postgres://orderly:secret@localhost:5432/orderly?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up
```

Seeded users: `jeremy@example.com` / `vijay@example.com`, password `password123`.

## Running tests

TODO

## Architecture

TODO — layering, request flow, async worker, webhook path.

## Data model

TODO — tables, keys, indexes and why.

## Trade-offs

Decisions made and what would change with more time.

- TODO

## Performance

`GET /orders` latency before and after Redis caching.

- TODO

## Assumptions

- TODO
