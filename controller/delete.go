package controller

import (
	"context"

	proto "github.com/open-Q/common/golang/proto/user"
	storageModel "github.com/open-Q/user/storage/model"
)

// Delete deletes an existing user.
func (s Service) Delete(ctx context.Context, req *proto.DeleteRequest, resp *proto.UserResponse) error {
	if err := s.userStorage.Delete(ctx, req.Id); err != nil {
		return err
	}
	users, err := s.userStorage.Find(ctx, storageModel.UserFindFilter{
		IDs: []string{req.Id},
	})
	if err != nil {
		return err
	}
	return newUserResponse(resp, &users[0])
}
