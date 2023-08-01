# retrier

[![Go Reference](https://pkg.go.dev/badge/github.com/craigpastro/retrier.svg)](https://pkg.go.dev/github.com/craigpastro/retrier)
[![Go Report Card](https://goreportcard.com/badge/github.com/craigpastro/retrier)](https://goreportcard.com/report/github.com/craigpastro/retrier)
[![CI](https://github.com/craigpastro/retrier/actions/workflows/ci.yaml/badge.svg)](https://github.com/craigpastro/retrier/actions/workflows/ci.yaml)
[![codecov](https://codecov.io/github/craigpastro/retrier/branch/main/graph/badge.svg?token=00AJODX77Z)](https://codecov.io/github/craigpastro/retrier)

A simple Go (Golang) library for retries with generics. Backoff is configurable
but the most useful is probably exponential backoff with full jitter.

See the [godoc](https://pkg.go.dev/github.com/craigpastro/retrier) for more
information.

## Usage

### Making an HTTP request

```go
body, err := retrier.DoWithData(func() ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}, retrier.NewExponentialBackoff())
```

### Connecting to a DB

```go
err := retrier.Do(func() error {
	if err = pool.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}, retrier.NewExponentialBackOff())
if err != nil {
	return nil, fmt.Errorf("error connecting to db: %w", err)
}
```

## Contributions

Contributions are welcome! Please create an issue for significant changes.

## Inspired By

- https://github.com/avast/retry-go
- https://github.com/cenkalti/backoff
