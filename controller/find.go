package controller

import (
	"context"

	proto "github.com/open-Q/common/golang/proto/user"
)

// Find finds users using filter.
func (s Service) Find(ctx context.Context, req *proto.FindFilter, resp proto.User_FindStream) error {
	// TODO: implement
	return nil
}
