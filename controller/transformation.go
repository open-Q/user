package controller

import (
	"github.com/golang/protobuf/ptypes"
	proto "github.com/open-Q/common/golang/proto/user"
	"github.com/open-Q/user/storage"
	"google.golang.org/protobuf/types/known/anypb"
)

type any struct {
	*anypb.Any
}

func newAny(value []byte) any {
	return any{
		Any: &anypb.Any{
			Value: value,
		},
	}
}

func newUserMeta(meta []*proto.UserMeta) map[string][]byte {
	res := make(map[string][]byte)
	for i := range meta {
		res[meta[i].Key] = meta[i].Value.GetValue()
	}
	return res
}

func newUserMetaProto(meta map[string][]byte) ([]*proto.UserMeta, error) {
	res := make([]*proto.UserMeta, 0, len(meta))
	for k, v := range meta {
		// skip nil values.
		if len(v) == 0 {
			continue
		}
		any, err := ptypes.MarshalAny(newAny(v))
		if err != nil {
			// skip value.
			continue
		}
		res = append(res, &proto.UserMeta{
			Key: k,
			Value: &anypb.Any{
				Value: any.Value,
			},
		})
	}
	return res, nil
}

func newUserResponse(resp *proto.UserResponse, user *storage.User) (err error) {
	resp.Id = user.ID
	resp.Status = proto.AccountStatus(proto.AccountStatus_value[user.Status])
	resp.Meta, err = newUserMetaProto(user.Meta)
	return
}
