package retrier

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

const DefaultAttempts = 100

var (
	errTimeoutExceeded = errors.New("timeout exceeded")
	errMaxAttempts     = errors.New("max attempts reached")
)

type RetryableFunc func() error
type RetryableFuncWithData[T any] func() (T, error)

type Config struct {
	Base       time.Duration
	Multiplier int
	// Retry will end after the last function call after Timeout.
	Timeout time.Duration
	// 0 means DefaultAttempts. Timeout takes precedence over Attempts.
	Attempts int
	// Uses "full-jitter" as described in https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/.
	Jitter bool
}

// NewDefault is exponential backoff with jitter and timeout of one minute.
func NewDefault() Config {
	return NewExponentialBackoff()
}

// NewExponentialBackoff with jitter and timeout of one minute.
func NewExponentialBackoff() Config {
	return Config{
		Base:       time.Second,
		Multiplier: 2,
		Timeout:    time.Minute,
		Jitter:     true,
	}
}

// NewConstantBackoff with user defined base duration and timeout of one minute.
func NewConstantBackoff(base time.Duration) Config {
	return Config{
		Base:       base,
		Multiplier: 1,
		Timeout:    time.Minute,
	}
}

// DoWithData will retry the RetryableFuncWithData in accordance with the given *Config.
func DoWithData[T any](f RetryableFuncWithData[T], cfg Config) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if cfg.Attempts == 0 {
		cfg.Attempts = DefaultAttempts
	}

	var emptyT T
	var lastErr error
	for i := 0; i < cfg.Attempts; i++ {
		t, err := f()
		if err == nil {
			return t, nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return emptyT, errors.Join(lastErr, errTimeoutExceeded)
		}

		sleep := cfg.Base.Milliseconds() * pow(cfg.Multiplier, i)
		if cfg.Jitter {
			sleep = rand.Int63n(sleep) + 1
		}

		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

	return emptyT, errors.Join(lastErr, errMaxAttempts)
}

// Do will retry the RetryableFunc in accordance with the given *Config.
func Do(f RetryableFunc, cfg Config) error {
	_, err := DoWithData(func() (any, error) {
		return nil, f()
	}, cfg)

	return err
}

func pow(x, y int) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}
