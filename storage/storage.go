package storage

import (
	"context"

	"github.com/open-Q/user/storage/model"
)

// User represents user's storage layer interface.
type User interface {
	Disconnect(ctx context.Context) error
	Add(ctx context.Context, user model.User) (*User, error)
	Delete(ctx context.Context, userID string) (*User, error)
	Update(ctx context.Context, user model.User) (*User, error)
	Find(ctx context.Context, filter model.UserFindFilter) ([]User, error)
}
