package controller

import (
	"github.com/micro/go-micro/v2"
	commonLog "github.com/open-Q/common/golang/log"
	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/dep"
	"github.com/open-Q/user/storage"
)

// Service represents service controller instance.
type Service struct {
	storage dep.Storage
	logger  *commonLog.Logger
}

// Config represents service configuration.
type Config struct {
	Storage dep.Storage
	Logger  *commonLog.Logger
	Micro   micro.Service
}

// New creates new service instance.
func New(cfg Config) (*Service, error) {
	srv := &Service{
		logger:  cfg.Logger,
		storage: cfg.Storage,
	}
	if err := proto.RegisterUserHandler(cfg.Micro.Server(), srv); err != nil {
		return nil, err
	}
	return srv, nil
}

func newUserResponse(resp *proto.UserResponse, user *storage.User) {
	resp.Email = user.Email
	if user.ID != nil {
		resp.Id = &proto.UserID{
			Id: *user.ID,
		}
	}
	resp.Status = proto.AccountStatus(proto.AccountStatus_value[user.Status])
}
