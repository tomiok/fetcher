## Fetcher service


This API returns the Last Traded Price of Bitcoin for the following currency pairs:

1. BTC/USD
2. BTC/CHF
3. BTC/EUR

### Requirements 
1. go 1.22
2. Docker (not mandatory)
3. [Mockgen](https://github.com/uber-go/mock)
4. Make (built-in for unix like OS)

## Make commands

### Tests
```shell
make tests
```

### Integration Tests
```shell
make tests-it
```

### Run
```shell
make run
```

### Run dockerized
```shell
make run-docker
```

---

### Run with Go (from the root)
```shell
go run cmd/api/main/main.go
```

### Run stand-alone docker
```shell
docker build -t fetcher .
docker run fetcher
```

### Test locally
(you can send any number of pairs)
(could be capitalized or not, is the same, but the format is with slash i.e: BTC/USD)
```shell
curl --location 'localhost:9000/api/v1/ltp?pairs=BTC%2FUSD%2CBTC%2FEUR%2CBTC%2FCHF'
```
or with one pair

```shell
curl --location 'localhost:9000/api/v1/ltp?pairs=btc%2Fusd'
```