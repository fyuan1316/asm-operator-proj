package manage

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func loopUntil(ctx context.Context, interval time.Duration, maxRetries int, f func(client.Client) (bool, error), client client.Client) error {
	count := 0
	for {
		if stop, err := f(client); err != nil {
			if count++; count > maxRetries {
				return err
			}
		} else if stop {
			break
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return fmt.Errorf("execute canceled")
		}
	}
	return nil
}
