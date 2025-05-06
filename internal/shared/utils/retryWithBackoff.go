package utils

import (
	"context"
	"errors"
	"log"
	"time"
)

func RetryWithBackoff(ctx context.Context, attempts int, fn func() error) error {
	if attempts <= 0 {
		return errors.New("attempts must be > 0")
	}

	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		// последняя попытка — возвращаем ошибку
		if i == attempts-1 {
			break
		}

		backoff := time.Duration(1<<i) * time.Second // 1s, 2s, 4s, 8s, ...
		log.Printf("retry %d/%d after %v: %v", i+1, attempts, backoff, err)

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
