# Reetro Application (Go)

Retrospective application built with Go, Postgres, and Redis. The app runs locally using Docker Compose, supports DB migrations, and exposes a REST API.

## Prerequisites

- Docker and Docker Compose
- Go (only required if you run Go tooling outside Docker)

## Environment

The app reads environment variables from the process. For local Docker runs, Docker Compose loads them from `.env`.

Required keys in `.env`:
- `APPLICATION_PORT`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_SSL_MODE`
- `REDIS_HOST`
- `REDIS_PORT`
- `JWT_SECRET_KEY`
- `EMAIL_ID`
- `EMAIL_PASSWORD`
- `EMAIL_HOST`
- `EMAIL_PORT`

## Services and ports

- API server: `APPLICATION_PORT` (default `8080`)
- Postgres: host `5400` → container `5432`
- Redis: host `6379` → container `6379`

## Quick start

These commands are kept as-is from the original README for consistency.

to build the application
- docker compose build

to run the application
- docker compose up

to stop the application
- docker compose down

to remove orphan containers of the application
- docker compose down --remove-orphans

## Development and testing

to inspect the container
- docker container inspect {container id}

to run the test cases
- docker compose run golang go test -cover -v ./...

## Makefile shortcuts

These wrap the Docker Compose commands:
- `make build`
- `make start` or `make up`
- `make down`
- `make test`
- `make migrate-up`
- `make migrate-down`
- `make makemigration name=init`

## Database migrations

to create Migrationion file
- docker compose run golang migrate create -ext=sql -dir={dir name or location} -seq {migration file name}
- example: docker compose run golang migrate create -ext=sql -dir=storages/migrations -seq init

to migrate Migration file up
- docker compose run golang migrate -path {full dir path or location} -database "postgres://{user}:{password}@{host}:5432/{db name}?sslmode={ssl_mode}" up
- example: docker compose run golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" up

to migrate Migration file down
- docker compose run golang migrate -path {full dir path or location} -database "postgres://{user}:{password}@{host}:5432/{db name}?sslmode={ssl_mode}" down
- example: docker compose run golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" down

## API routes (high level)

Public routes:
- `/`
- `/about/`
- `/login/`
- `/signup/`
- `/forgot_password/`
- `/reset_password/`
- `/clear_redis/`

Protected routes (JWT required):
- `/users/`
- `/users/{id}/`
- `/boards/`
- `/boards/{id}/`
- `/feedbacks/`
- `/feedbacks/{id}/`

## Access Postgres container

to access postgres database
- postgres container should be in running state
- docker exec -it {image name} /bin/sh
- example: docker exec -it reetro_postgres /bin/sh
- after entering into postgres docker container run command :-> psql --username {db username}
- example: psql --username postgres
- \l command for list of databases
- then enter into our application database
- command: \c {database name};
- example: \c reetro_app;
- command: \dt -> list of tables
- command: \d {table name}; -> will show schema of database table
- example: \d users;

- for quit or get back to the main option:-> \q
- to come out of postgres docker image:-> exit

## Go module maintenance

to remove unwanted packages
- go mod tidy

## CI

CI runs on GitHub Actions and executes Go tests with coverage. The workflow uses the same Go version as `go.mod` and prints a coverage summary in the logs.