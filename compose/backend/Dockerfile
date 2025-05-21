FROM golang:1.23-alpine AS base

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

FROM base AS dev

EXPOSE 8080

CMD ["air"]
