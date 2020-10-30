package controller

import (
	"errors"
	"testing"

	"github.com/micro/go-micro/v2"
	commonLog "github.com/open-Q/common/golang/log"
	"github.com/open-Q/user/mocks"
	"github.com/stretchr/testify/require"
)

var errMock = errors.New("error")

func Test_New(t *testing.T) {
	s, err := New(Config{
		Storage: &mocks.Storage{},
		Micro:   micro.NewService(),
		Logger:  &commonLog.Logger{},
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, s.logger)
	require.NotNil(t, s.storage)
}
