package controller

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/storage"
	"github.com/stretchr/testify/require"
)

func Test_newAny(t *testing.T) {
	v := []byte("hello")
	any := newAny(v)
	require.NotNil(t, any)
	require.Equal(t, v, any.Value)
}

func Test_newUserMeta(t *testing.T) {
	meta := []*proto.UserMeta{
		{
			Key:   "key1",
			Value: newAny([]byte("hello")).Any,
		},
		{
			Key:   "key2",
			Value: newAny([]byte("world")).Any,
		},
	}
	userMeta := newUserMeta(meta)
	require.NotNil(t, userMeta)
	for i := range meta {
		v, ok := userMeta[meta[i].Key]
		require.True(t, ok)
		require.True(t, bytes.Equal(v, meta[i].Value.GetValue()))
	}
}

func Test_newUserMetaProto(t *testing.T) {
	meta := map[string][]byte{
		"hello": nil,
		"key1":  []byte("value"),
		"key2":  {1, 2, 3},
	}
	res, err := newUserMetaProto(meta)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, len(meta)-1, len(res))
	for k, v := range meta {
		if len(v) == 0 {
			continue
		}
		var found bool
		for i := range res {
			if res[i].Key == k {
				any, err := ptypes.MarshalAny(newAny(v))
				require.NoError(t, err)
				require.NotNil(t, any)
				require.True(t, bytes.Equal(res[i].Value.GetValue(), any.Value))
				found = true
				break
			}
		}
		require.True(t, found)
	}
}

func Test_newUserResponse(t *testing.T) {
	user := storage.User{
		ID:     "1",
		Status: "ACCOUNT_STATUS_ACTIVE",
		Meta: map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		},
	}
	var resp proto.UserResponse
	err := newUserResponse(&resp, &user)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "1", resp.Id)
	require.Equal(t, proto.AccountStatus_ACCOUNT_STATUS_ACTIVE, resp.Status)
	require.NotNil(t, resp.Meta)
	require.Equal(t, len(user.Meta), len(resp.Meta))
	for k, v := range user.Meta {
		var found bool
		for i := range resp.Meta {
			if resp.Meta[i].Key == k {
				any, err := ptypes.MarshalAny(newAny(v))
				require.NoError(t, err)
				require.NotNil(t, any)
				require.True(t, bytes.Equal(any.Value, resp.Meta[i].Value.GetValue()))
				found = true
				break
			}
		}
		require.True(t, found)
	}
}
