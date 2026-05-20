# Supplier service

## Run

```bash
export CONFIG_PATH=configs/supplier.env

go run ./services/supplier-service/cmd/supplier-service
```

## Run via Makefile

```bash
make supplier-run
```

## Docker

Build and run only service container:

```bash
make supplier-docker-build
make supplier-docker-run
```

Run service + PostgreSQL (recommended for local docker run):

```bash
make supplier-compose-up
make supplier-compose-migrate
make supplier-compose-logs
```

Stop and remove compose stack:

```bash
make supplier-compose-down
```

## Tests

```bash
go test ./services/supplier-service/...
```


