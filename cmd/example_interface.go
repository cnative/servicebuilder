package cmd

import (
	"context"
	"io"
)

//Store example interface
type Store interface {
	Serve(ctx context.Context) error
	io.Closer

	Create(ctx context.Context, c *string) (*string, error)
	Get(ctx context.Context, id string) (string, error)
	Filter(ctx context.Context, filterReq string, opts ...io.Closer) ([]*string, error)
}
