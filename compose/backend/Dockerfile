ARG GO_VERSION=1.25.7
FROM golang:${GO_VERSION}-alpine AS base

ARG MIGRATE_VERSION=v4.19.1
ARG AIR_VERSION=v1.64.5

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@${MIGRATE_VERSION}

RUN go install github.com/air-verse/air@${AIR_VERSION}

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

FROM base AS dev

EXPOSE 8080

CMD ["air"]
