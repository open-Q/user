package controller

import (
	"context"

	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/storage"
)

// Create creates new user.
func (s Service) Create(ctx context.Context, req *proto.UserRequest, resp *proto.UserResponse) error {
	user := storage.User{
		Email:  req.GetEmail(),
		Status: proto.AccountStatus_ACCOUNT_STATUS_ACTIVE.String(),
	}
	createdUser, err := s.storage.Add(ctx, user)
	if err != nil {
		return err
	}
	newUserResponse(resp, createdUser)
	return nil
}
