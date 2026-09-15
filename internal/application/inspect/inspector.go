package inspect

import "context"

type Inspector interface {
	Inspect(ctx context.Context, url string) (Result, error)
}
