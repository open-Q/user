package controller

import (
	"github.com/micro/go-micro/v2"
	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/dep"
)

// Service represents service controller instance.
type Service struct {
	storage dep.Storage
}

// New creates new service instance.
func New(s micro.Service) (*Service, error) {
	srv := &Service{}
	if err := proto.RegisterUserHandler(s.Server(), srv); err != nil {
		return nil, err
	}
	return srv, nil
}
