This application is used for retrospective.

to build the application
- docker compose build

to run the application
- docker compose up

to stop the application
- docker compose down

to remove orphan containers of the application
- docker compose down --remove-orphans

to inspect the container
- docker container inspect {container id}

to run the test cases
- docker compose run golang go test -cover -v ./...

to create Migrationion file
- docker compose run golang migrate create -ext=sql -dir={dir name or location} -seq {migration file name}
- example: docker compose run golang migrate create -ext=sql -dir=storages/migrations -seq init

to migrate Migration file up
- docker compose run golang migrate -path {full dir path or location} -database "postgres://{user}:{password}@{host}:5432/{db name}?sslmode={ssl_mode}" up
- example: docker compose run golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" up

to migrate Migration file down
- docker compose run golang migrate -path {full dir path or location} -database "postgres://{user}:{password}@{host}:5432/{db name}?sslmode={ssl_mode}" down
- example: docker compose run golang migrate -path /app/storages/migrations -database "postgres://postgres:postgres@postgres:5432/reetro_app?sslmode=disable" down

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

to remove unwanted packages
- go mod tidy