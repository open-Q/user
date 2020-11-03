package controller

import (
	"context"

	proto "github.com/open-Q/common/golang/proto/user"
	storageModel "github.com/open-Q/user/storage/model"
)

// Create creates new user.
func (s Service) Create(ctx context.Context, req *proto.CreateRequest, resp *proto.UserResponse) error {
	user := storageModel.User{
		Status: proto.AccountStatus_name[int32(proto.AccountStatus_ACCOUNT_STATUS_ACTIVE)],
		Meta:   newUserMeta(req.Meta),
	}
	createdUser, err := s.userStorage.Add(ctx, user)
	if err != nil {
		return err
	}
	return newUserResponse(resp, createdUser)
}
