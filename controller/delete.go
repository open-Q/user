package controller

import (
	"context"

	proto "github.com/open-Q/common/golang/proto/user"
)

// Delete deletes existing user.
func (s Service) Delete(ctx context.Context, req *proto.UserID, resp *proto.UserID) error {
	// TODO: implement
	return nil
}
