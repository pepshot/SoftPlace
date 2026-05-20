# syntax=docker/dockerfile:1

FROM golang:1.25.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/customer-service ./services/customer-service/cmd/customer-service

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/customer-service /app/customer-service
COPY configs/customer.docker.env /app/configs/customer.docker.env

ENV CONFIG_PATH=/app/configs/customer.docker.env

EXPOSE 8442

USER app

CMD ["/app/customer-service"]

