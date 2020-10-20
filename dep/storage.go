package dep

import (
	"context"

	"github.com/open-Q/user/storage"
)

// Storage represents storage layer interface.
type Storage interface {
	Disconnect(ctx context.Context) error
	Add(ctx context.Context, user storage.User) (*storage.User, error)
}
