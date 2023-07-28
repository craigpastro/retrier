# retrier

A simple, no dependency, Go (Golang) library for retries.

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
