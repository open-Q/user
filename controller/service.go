package controller

import (
	commonLog "github.com/open-Q/common/golang/log"
	"github.com/open-Q/user/storage"
)

// Service represents service controller instance.
type Service struct {
	userStorage storage.User
	logger      *commonLog.Logger
}

// Config represents service configuration.
type Config struct {
	UserStorage storage.User
	Logger      *commonLog.Logger
}

// New creates new service instance.
func New(cfg Config) Service {
	return Service{
		logger:      cfg.Logger,
		userStorage: cfg.UserStorage,
	}
}
