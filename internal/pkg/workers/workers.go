package workers

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func Timeout(ctx context.Context, d time.Duration, fn func(ctx context.Context) error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(d):
				if err := fn(ctx); err != nil {
					log.Fatal().
						Err(err).
						Msg("worker timeout error")

					return
				}
			}
		}
	}()
}
