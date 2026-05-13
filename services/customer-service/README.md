# Customer service

## Run

```bash
export CONFIG_PATH=configs/customer.env

go run ./services/customer-service/cmd/customer-service
```

## Swagger

Open `http://localhost:8442/swagger/index.html` after the service starts.

## Tests

```bash
go test ./services/customer-service/...
```
