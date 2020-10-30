package controller

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2"
	"github.com/open-Q/common/golang/proto/user"
	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/mocks"
	"github.com/open-Q/user/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	service, err := New(Config{
		Micro: micro.NewService(),
	})
	require.NoError(t, err)
	require.NotNil(t, service)
	t.Run("save to storage error", func(t *testing.T) {
		st := new(mocks.Storage)
		defer st.AssertExpectations(t)
		service.storage = st
		st.On("Add", mock.Anything, mock.Anything).Return(nil, errMock)
		err = service.Create(context.Background(), &user.CreateRequest{}, &user.UserResponse{})
		require.Error(t, err)
		require.EqualError(t, err, errMock.Error())
	})
	t.Run("all ok", func(t *testing.T) {
		st := new(mocks.Storage)
		defer st.AssertExpectations(t)
		service.storage = st
		req := &user.CreateRequest{
			Meta: []*proto.UserMeta{
				{
					Key:   "key1",
					Value: newAny([]byte("value1")).Any,
				},
				{
					Key:   "key2",
					Value: newAny([]byte{1, 2, 3}).Any,
				},
			},
		}
		storageResponse := storage.User{
			ID:     "1",
			Status: proto.AccountStatus_name[int32(proto.AccountStatus_ACCOUNT_STATUS_ACTIVE)],
			Meta:   newUserMeta(req.Meta),
		}
		st.On("Add", mock.Anything, storage.User{
			Status: proto.AccountStatus_name[int32(proto.AccountStatus_ACCOUNT_STATUS_ACTIVE)],
			Meta:   newUserMeta(req.Meta),
		}).Return(&storageResponse, nil)
		var resp user.UserResponse
		err = service.Create(context.Background(), req, &resp)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, storageResponse.ID, resp.Id)
		require.Equal(t, proto.AccountStatus(proto.AccountStatus_value[storageResponse.Status]), resp.Status)
		require.NotNil(t, resp.Meta)
		require.Equal(t, len(req.Meta), len(resp.Meta))
		for i := range req.Meta {
			var found bool
			for j := range resp.Meta {
				if resp.Meta[j].Key == req.Meta[i].Key {
					any, err := ptypes.MarshalAny(newAny(req.Meta[i].Value.GetValue()))
					require.NoError(t, err)
					require.True(t, bytes.Equal(resp.Meta[j].Value.GetValue(), any.Value))
					found = true
					break
				}
			}
			require.True(t, found)
		}
	})
}
