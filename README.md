# REST Pricing Server

## Getting started

Start a Postgres database with Docker:

```
docker run \
  --rm \
  --name postgres \
  --publish 5432:5432 \
  --env POSTGRES_PASSWORD=password \
  --detach \
  docker.io/postgres:15.3-alpine
```

Run the application:

```
go run ./cmd/server/main.go
```

Query for pricing:

```
curl "http://127.0.0.1:8080/api/v1/pricing?date=1&productID=2&stringID=3"

# response:
# TODO
```

## Contributing

Run tests:

```
go test -race -v ./...
```

Run tests without testcontainers (Docker):

```
go test -short -race -v ./...
```
