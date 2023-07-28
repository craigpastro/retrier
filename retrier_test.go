package retrier

import (
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("error")

func TestRetry(t *testing.T) {
	cfg := DefaultConfig()

	var attempts int
	err := Do(func() error {
		attempts += 1
		if attempts == 3 {
			return nil
		}
		return errTest
	}, cfg)

	if err != nil {
		t.Error(err)
	}

	if attempts != 3 {
		t.Error("number of attempts is incorrect:", attempts)
	}
}

func TestZeroTimeoutWillTryOnce(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Timeout = 0

	var attempts int
	err := Do(func() error {
		attempts += 1
		return errTest

	}, cfg)

	if !errors.Is(err, errTest) {
		t.Error(err)
	}

	if attempts != 1 {
		t.Error("number of attempts is incorrect:", attempts)
	}
}

func TestConstantBackoff(t *testing.T) {
	cfg := ConstantBackoff(time.Millisecond)
	cfg.Attempts = 10
	leastDuration := 10 * time.Millisecond
	mostDuration := 15 * time.Millisecond
	now := time.Now()

	var attempts int
	err := Do(func() error {
		attempts += 1
		return errTest
	}, cfg)

	if !errors.Is(err, errTest) {
		t.Error(err)
	}

	if attempts != 10 {
		t.Error("number of attempts is incorrect:", attempts)
	}

	since := time.Since(now)

	if since < leastDuration {
		t.Errorf("backoff duration is wrong. expected at least %v, got: %v", leastDuration, since)
	}

	if since > mostDuration {
		t.Errorf("backoff duration is wrong. expected at most %v, got: %v", mostDuration, since)
	}
}

func TestExponentialBackoff(t *testing.T) {
	cfg := Config{
		Base:       time.Millisecond,
		Multiplier: 2,
		Timeout:    time.Minute,
		Attempts:   4, // should take about 15ms
	}
	leastDuration := 15 * time.Millisecond
	mostDuration := 20 * time.Millisecond
	now := time.Now()

	var attempts int
	err := Do(func() error {
		attempts += 1
		return errTest
	}, cfg)

	if !errors.Is(err, errTest) {
		t.Error(err)
	}

	if attempts != 4 {
		t.Error("number of attempts is incorrect:", attempts)
	}

	since := time.Since(now)

	if since < leastDuration {
		t.Errorf("backoff duration is wrong. expected at least %v, got: %v", leastDuration, since)
	}

	if since > mostDuration {
		t.Errorf("backoff duration is wrong. expected at most %v, got: %v", mostDuration, since)
	}
}
