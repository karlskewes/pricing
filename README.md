# REST Pricing Server

## Getting started

Run the application:

```
go run ./cmd/server/main.go
```

Query for brands:

```
# Default seeded brand
curl localhost:8080/api/v1/brands?name=EXAMPLE
{"id":1,"name":"EXAMPLE"}

# Non-existant brand
curl -I localhost:8080/api/v1/brands?name=NOTEXIST
HTTP/1.1 404 Not Found
```

Query for pricing, note time is in [RFC3339](https://en.wikipedia.org/wiki/ISO_8601#RFCs):

```
curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T10:00:00.00Z&string_id=test_1' | jq -r
{
  "brand_id": 1,
  "product_id": 35455,
  "price": "35.50",
  "curr": "EUR",
  "start_date": "2020-06-14 00:00:00 +0000 UTC",
  "end_date": "2020-12-31 23:59:59 +0000 UTC",
  "string_id": "test_1"
}
```

All pricing queries:

```
curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T10:00:00.00Z&string_id=test_1' | jq -r

curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T16:00:00.00Z&string_id=test_1' | jq -r

curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T21:00:00.00Z&string_id=test_1' | jq -r

curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-15T10:00:00.00Z&string_id=test_1' | jq -r

curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-16T21:00:00.00Z&string_id=test_1' | jq -r
```

## Use Postgres repository

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
go run ./cmd/server/main.go -enable-postgres=true
```

## Tests

Run tests:

```
go test -race -v ./...
```

Run tests without testcontainers (Docker):

```
go test -short -race -v ./...
```

Run end to end test with provided date times:

```
go test -race -v ./... -run TestRun

=== RUN   TestRun
=== RUN   TestRun/Test_1
    server_test.go:142: test_1 - {1 35455 35.50 EUR 2020-06-14 00:00:00 +0000 UTC 2020-12-31 23:59:59 +0000 UTC test_1}
=== RUN   TestRun/Test_2
    server_test.go:142: test_2 - {1 35455 25.45 EUR 2020-06-14 15:00:00 +0000 UTC 2020-06-14 18:30:00 +0000 UTC test_2}
=== RUN   TestRun/Test_3
    server_test.go:142: test_3 - {1 35455 35.50 EUR 2020-06-14 00:00:00 +0000 UTC 2020-12-31 23:59:59 +0000 UTC test_3}
=== RUN   TestRun/Test_4
    server_test.go:142: test_4 - {1 35455 30.50 EUR 2020-06-15 00:00:00 +0000 UTC 2020-06-15 11:00:00 +0000 UTC test_4}
=== RUN   TestRun/Test_5
    server_test.go:142: test_5 - {1 35455 38.95 EUR 2020-06-15 16:00:00 +0000 UTC 2020-12-31 23:59:59 +0000 UTC test_5}
--- PASS: TestRun (0.00s)
    --- PASS: TestRun/Test_1 (0.00s)
    --- PASS: TestRun/Test_2 (0.00s)
    --- PASS: TestRun/Test_3 (0.00s)
    --- PASS: TestRun/Test_4 (0.00s)
    --- PASS: TestRun/Test_5 (0.00s)
```

## Shortcuts

Many!

- doc comments
- not using https://github.com/Rhymond/go-money
- table tests for all situations
- all API's
- context cancellations
- no logger like zerolog, zap, logrus, etc
- no Prometheus metrics
- no OTEL
- hardening
- `-version` command based on git sha
- OpenAPI yaml

## Results

Same for in memory and postgres backend repositories.

```
$ curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T10:00:00.00Z&string_id=test_1' | jq -r
{
  "brand_id": 1,
  "product_id": 35455,
  "price": "35.50",
  "curr": "EUR",
  "start_date": "2020-06-14 00:00:00 +0000 UTC",
  "end_date": "2020-12-31 23:59:59 +0000 UTC",
  "string_id": "test_1"
}

$ curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T16:00:00.00Z&string_id=test_2' | jq -r
{
  "brand_id": 1,
  "product_id": 35455,
  "price": "25.45",
  "curr": "EUR",
  "start_date": "2020-06-14 15:00:00 +0000 UTC",
  "end_date": "2020-06-14 18:30:00 +0000 UTC",
  "string_id": "test_2"
}

$ curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-14T21:00:00.00Z&string_id=test_3' | jq -r
{
  "brand_id": 1,
  "product_id": 35455,
  "price": "35.50",
  "curr": "EUR",
  "start_date": "2020-06-14 00:00:00 +0000 UTC",
  "end_date": "2020-12-31 23:59:59 +0000 UTC",
  "string_id": "test_3"
}

$ curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-15T10:00:00.00Z&string_id=test_4' | jq -r
{
  "brand_id": 1,
  "product_id": 35455,
  "price": "30.50",
  "curr": "EUR",
  "start_date": "2020-06-15 00:00:00 +0000 UTC",
  "end_date": "2020-06-15 11:00:00 +0000 UTC",
  "string_id": "test_4"
}

$ curl -s 'localhost:8080/api/v1/prices?brand_id=1&product_id=35455&date=2020-06-16T21:00:00.00Z&string_id=test_5' | jq -r
{
  "brand_id": 1,
  "product_id": 35455,
  "price": "38.95",
  "curr": "EUR",
  "start_date": "2020-06-15 16:00:00 +0000 UTC",
  "end_date": "2020-12-31 23:59:59 +0000 UTC",
  "string_id": "test_5"
}
```
