# Customer service

## Run

```bash
export CONFIG_PATH=configs/customer.env

go run ./services/customer-service/cmd/customer-service
```

## Run via Makefile

```bash
make customer-run
```

## Swagger

Open `http://localhost:8442/swagger/index.html` after the service starts.

## Docker

Build and run only service container:

```bash
make customer-docker-build
make customer-docker-run
```

Run service + PostgreSQL (recommended for local docker run):

```bash
make customer-compose-up
make customer-compose-migrate
make customer-compose-logs
```

Stop and remove compose stack:

```bash
make customer-compose-down
```

## Tests

```bash
go test ./services/customer-service/...
```
