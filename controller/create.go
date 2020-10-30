package controller

import (
	"context"

	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/storage"
)

// Create creates new user.
func (s Service) Create(ctx context.Context, req *proto.CreateRequest, resp *proto.UserResponse) error {
	user := storage.User{
		Status: proto.AccountStatus_ACCOUNT_STATUS_ACTIVE.String(),
		Meta:   newUserMeta(req.Meta),
	}
	createdUser, err := s.storage.Add(ctx, user)
	if err != nil {
		return err
	}
	return newUserResponse(resp, createdUser)
}
