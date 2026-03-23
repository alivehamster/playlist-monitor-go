package utils

import (
	"context"
	"log"
	"time"
)

func StartDaily(ctx context.Context, job func() error) {
	go func() {
		for {
			if err := job(); err != nil {
				log.Printf("daily job error: %v", err)
			}

			select {
			case <-time.After(24 * time.Hour):
			case <-ctx.Done():
				return
			}
		}
	}()
}
