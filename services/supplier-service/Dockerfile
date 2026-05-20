# syntax=docker/dockerfile:1

FROM golang:1.25.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/supplier-service ./services/supplier-service/cmd/supplier-service

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/supplier-service /app/supplier-service
COPY configs/supplier.docker.env /app/configs/supplier.docker.env

ENV CONFIG_PATH=/app/configs/supplier.docker.env

EXPOSE 8443

USER app

CMD ["/app/supplier-service"]


