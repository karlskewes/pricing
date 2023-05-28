# REST Pricing Server

## Getting started

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
