package retrier

import (
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("error")

func TestRetry(t *testing.T) {
	cfg := NewDefault()
	cfg.Base = time.Millisecond

	var i int
	attempts, err := DoWithData(func() (int, error) {
		i += 1
		if i == 3 {
			return i, nil
		}
		return 0, errTest
	}, cfg)

	if err != nil {
		t.Error(err)
	}

	if attempts != 3 {
		t.Error("number of attempts is incorrect:", attempts)
	}
}

func TestZeroTimeoutWillTryOnce(t *testing.T) {
	cfg := NewDefault()
	cfg.Timeout = 0

	var i int
	err := Do(func() error {
		i += 1
		return errTest
	}, cfg)

	if !errors.Is(err, errTimeoutExceeded) {
		t.Error(err)
	}

	if i != 1 {
		t.Error("number of attempts is incorrect:", i)
	}
}

func TestConstantBackoffAndMaxAttempts(t *testing.T) {
	cfg := NewConstantBackoff(time.Millisecond)
	cfg.Attempts = 10
	leastDuration := 10 * time.Millisecond
	mostDuration := 15 * time.Millisecond
	now := time.Now()

	var attempts int
	err := Do(func() error {
		attempts += 1
		return errTest
	}, cfg)

	if !errors.Is(err, errMaxAttempts) && !errors.Is(err, errTest) {
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

func TestExponentialBackoffAndMaxAttempts(t *testing.T) {
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

	if !errors.Is(err, errMaxAttempts) && !errors.Is(err, errTest) {
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
