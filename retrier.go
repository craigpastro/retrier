package retrier

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

const (
	DefaultAttempts = 100
	DefaultTimeout  = time.Minute
)

type RetryableFunc func() error

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

func DefaultConfig() Config {
	return ExponentialBackoff()
}

// ExponentialBackoff with jitter and timeout of one minute.
func ExponentialBackoff() Config {
	return Config{
		Base:       time.Second,
		Multiplier: 2,
		Timeout:    time.Minute,
		Jitter:     true,
	}
}

// ConstantBackoff with user defined base duration and timeout of one minute.
func ConstantBackoff(base time.Duration) Config {
	return Config{
		Base:       base,
		Multiplier: 1,
		Timeout:    time.Minute,
	}
}

// Do will retry the RetryableFunc in accordance with the given *Config.
func Do(f RetryableFunc, cfg Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if cfg.Attempts == 0 {
		cfg.Attempts = DefaultAttempts
	}

	var lastErr error
	for i := 0; i < cfg.Attempts; i++ {
		err := f()
		if err == nil {
			return nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return errors.Join(lastErr, errors.New("timeout exceeded"))
		}

		sleep := cfg.Base.Milliseconds() * pow(cfg.Multiplier, i)
		if cfg.Jitter {
			sleep = rand.Int63n(sleep) + 1
		}

		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

	return errors.Join(lastErr, errors.New("max attempts reached"))
}

func pow(x, y int) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}
