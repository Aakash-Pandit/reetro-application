build:
	docker compose build

stop:
	docker compose stop

up:
	make down
	make migrate-up
	docker compose up --remove-orphans

down:
	docker compose down --remove-orphans

remove-containers:
	@if [ -n "$$(docker ps -a)" ]; then docker rm $$(docker ps -a); else echo "No containers to remove."; fi

test:
	make down
	docker compose run --remove-orphans golang go test -v -cover ./...

clear-cache:
	docker compose run --remove-orphans golang go clean -cache -modcache

migrate-up:
	docker compose run --remove-orphans golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" up

migrate-down:
	docker compose run --remove-orphans golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" down

makemigration:
	docker compose run golang migrate create -ext=sql -dir=storages/migrations -seq ${name}

prod-build:
	docker-compose build

prod-stop:
	docker-compose stop

prod-up:
	make prod-stop
	make prod-down
	make prod-migrate-up
	docker-compose up --remove-orphans -d

prod-down:
	docker-compose down --remove-orphans

prod-remove-containers:
	@if [ -n "$$(docker ps -a)" ]; then docker rm $$(docker ps -a); else echo "No containers to remove."; fi

prod-test:
	make down
	docker-compose run --remove-orphans golang go test -v -cover ./...

prod-clear-cache:
	docker-compose run --remove-orphans golang go clean -cache -modcache

prod-migrate-up:
	docker-compose run --remove-orphans golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" up

prod-migrate-down:
	docker-compose run --remove-orphans golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" down

prod-makemigration:
	docker-compose run golang migrate create -ext=sql -dir=storages/migrations -seq ${name}