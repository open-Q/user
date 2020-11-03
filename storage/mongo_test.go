package storage

import (
	"context"
	"testing"

	commonErrors "github.com/open-Q/common/golang/errors"
	"github.com/open-Q/user/storage/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const testConnection = "mongodb://127.0.0.1:27017"

func Test_NewMongoStorage(t *testing.T) {
	t.Run("create connection error", func(t *testing.T) {
		_, err := NewMongoStorage(context.Background(), "invalid", "test-db")
		require.Error(t, err)
		require.EqualError(t, err, "could not create mongo client: error parsing uri: scheme must be \"mongodb\" or \"mongodb+srv\"")
	})
	t.Run("all ok", func(t *testing.T) {
		st, err := NewMongoStorage(context.Background(), testConnection, "test-db")
		require.NoError(t, err)
		require.NotNil(t, st)
		defer clearMongoStorage(t, st)
		require.NotNil(t, st.userCollection)
		require.NotNil(t, st.db)
	})
}

func TestMongoStorage_Disconnect(t *testing.T) {
	st, err := NewMongoStorage(context.Background(), testConnection, "test-db")
	require.NoError(t, err)
	err = st.Disconnect(context.Background())
	require.NoError(t, err)
}

func Test_NewMongoUser(t *testing.T) {
	t.Run("parse ID error", func(t *testing.T) {
		id := "invalid"
		user := model.User{
			ID:     id,
			Status: "some status",
		}
		_, err := NewMongoUser(user)
		require.Error(t, err)
		require.EqualError(t, err, "encoding/hex: invalid byte: U+0069 'i'")
	})
	t.Run("all ok", func(t *testing.T) {
		id := primitive.NewObjectID()
		idHex := id.Hex()
		user := model.User{
			ID:     idHex,
			Status: "some status",
			Meta: map[string]interface{}{
				"hello": "world",
				"key":   []int{1, 2, 3},
			},
		}
		res, err := NewMongoUser(user)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, user.ID, res.ID.Hex())
		require.Equal(t, user.Status, res.Status)
		require.NotNil(t, res.Meta)
		for k, v := range user.Meta {
			value, ok := res.Meta[k]
			require.True(t, ok)
			require.Equal(t, v, value)
		}
	})
}

func TestMongoUser_ToUser(t *testing.T) {
	id := primitive.NewObjectID()
	idHex := id.Hex()
	userMongo := MongoUser{
		ID:     id,
		Status: "some status",
		Meta: map[string]interface{}{
			"hello": "world",
			"key":   []int{1, 2, 3},
		},
	}
	res := userMongo.ToUser()
	require.NotNil(t, res)
	require.Equal(t, model.User{
		ID: idHex,
		Meta: map[string]interface{}{
			"hello": "world",
			"key":   []int{1, 2, 3},
		},
		Status: userMongo.Status,
	}, *res)
}

func TestMongoStorage_Add(t *testing.T) {
	t.Run("convertation error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := "invalid"
		_, err := st.Add(context.Background(), model.User{
			ID: id,
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageConvert))
	})
	t.Run("insert error (duplicate entry id)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := primitive.NewObjectID().Hex()
		_, err := st.Add(context.Background(), model.User{
			ID: id,
		})
		require.NoError(t, err)
		_, err = st.Add(context.Background(), model.User{
			ID: id,
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageInsert))
		require.Contains(t, err.Error(), "duplicate key error collection")
	})
	t.Run("all ok", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		user, err := st.Add(context.Background(), model.User{
			Status: "some status",
			Meta: map[string]interface{}{
				"hello": "world",
				"key":   []interface{}{int32(1), int32(2), int32(3)},
				"slice": []interface{}{
					[]interface{}{"1", "2", "3"},
					[]interface{}{int32(1), int32(2), int32(3)},
					[]interface{}{"1", "2", []interface{}{"1", "2", "3"}},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, user)
		var fUser MongoUser
		err = st.userCollection.FindOne(context.Background(), bson.M{}).Decode(&fUser)
		require.NoError(t, err)
		require.NotNil(t, fUser)
		require.Equal(t, *user, *fUser.ToUser())
	})
}

func TestMongoStorage_Delete(t *testing.T) {
	t.Run("convertation error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := "invalid"
		err := st.Delete(context.Background(), id)
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageConvert))
	})
	t.Run("nothing to delete error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := primitive.NewObjectID()
		err := st.Delete(context.Background(), id.Hex())
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageDelete))
		require.Contains(t, err.Error(), "user not found")
	})
	t.Run("all ok", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := []MongoUser{
			{
				ID: primitive.NewObjectID(),
			},
			{
				ID: primitive.NewObjectID(),
			},
		}
		docs := make([]interface{}, len(users))
		for i := range users {
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		err = st.Delete(ctx, users[0].ID.Hex())
		require.NoError(t, err)
		cur, err := st.userCollection.Find(ctx, bson.M{})
		require.NoError(t, err)
		defer closeCursor(ctx, cur)
		var results []MongoUser
		err = cur.All(ctx, &results)
		require.NoError(t, err)
		require.Equal(t, len(users)-1, len(results))
		require.Equal(t, users[1].ID, results[0].ID)
	})
}

func TestMongoStorage_Update(t *testing.T) {
	t.Run("convertation error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := "invalid"
		_, err := st.Update(context.Background(), model.User{
			ID: id,
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageConvert))
	})
	t.Run("nothing to update error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := primitive.NewObjectID()
		_, err := st.Update(context.Background(), model.User{
			ID: id.Hex(),
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageUpdate))
		require.Contains(t, err.Error(), "user not found")
	})
	t.Run("all ok", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := []MongoUser{
			{
				ID:     primitive.NewObjectID(),
				Status: "some status",
				Meta: map[string]interface{}{
					"key1": "world",
					"key2": []interface{}{int32(1), int32(2), int32(3)},
				},
			},
			{
				ID:     primitive.NewObjectID(),
				Status: "some status",
			},
		}
		docs := make([]interface{}, len(users))
		for i := range users {
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		userToUpdate := model.User{
			ID:     users[0].ID.Hex(),
			Status: "new status",
			Meta: map[string]interface{}{
				"key1": "world",
				"key3": []interface{}{int32(1), int32(2), int32(3)},
				"key4": map[string]interface{}{
					"key4_1": "value",
					"key4_2": []interface{}{"1", "2", "3"},
					"key4_3": map[string]interface{}{
						"key4_3_1": "hello",
						"key4_3_2": []interface{}{int32(1), int32(2), int32(3)},
					},
				},
			},
		}
		res, err := st.Update(ctx, userToUpdate)
		require.NoError(t, err)
		require.NotNil(t, res)
		var user MongoUser
		err = st.userCollection.FindOne(ctx, bson.M{}).Decode(&user)
		require.NoError(t, err)
		require.Equal(t, userToUpdate, *user.ToUser())
	})
}

func TestMongoStorage_Find(t *testing.T) {
	t.Run("create filter error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		_, err := st.Find(context.Background(), model.UserFindFilter{
			IDs: []string{"invalid"},
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, commonErrors.ErrStorageConvert))
	})
	t.Run("all ok (only ids in the filter)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := []MongoUser{
			{
				ID: primitive.NewObjectID(),
			},
			{
				ID: primitive.NewObjectID(),
			},
		}
		docs := make([]interface{}, len(users))
		for i := range users {
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			IDs: []string{users[0].ID.Hex()},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(resp))
		require.Equal(t, users[0].ID.Hex(), resp[0].ID)
	})
	t.Run("all ok (only statuses in the filter)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := []MongoUser{
			{
				ID:     primitive.NewObjectID(),
				Status: "some status",
			},
			{
				ID:     primitive.NewObjectID(),
				Status: "some status 2",
			},
			{
				ID:     primitive.NewObjectID(),
				Status: "some status 3",
			},
		}
		docs := make([]interface{}, len(users))
		for i := range users {
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			Statuses: []string{users[0].Status, users[2].Status},
		})
		require.NoError(t, err)
		require.Equal(t, 2, len(resp))
		userIDs := []string{users[0].ID.Hex(), users[2].ID.Hex()}
		for i := range resp {
			require.Contains(t, userIDs, resp[i].ID)
		}
	})
	t.Run("all ok (only meta in the filter)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := []MongoUser{
			{
				ID: primitive.NewObjectID(),
				Meta: map[string]interface{}{
					"email": "test@gmail.com",
					"key":   []int{1, 2, 3},
				},
			},
			{
				ID: primitive.NewObjectID(),
				Meta: map[string]interface{}{
					"key": []string{"1", "2", "3"},
				},
			},
			{
				ID: primitive.NewObjectID(),
				Meta: map[string]interface{}{
					"email": "test2@gmail.com",
					"key":   "some key",
				},
			},
		}
		docs := make([]interface{}, len(users))
		for i := range users {
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		// check all with not empty email value.
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			MetaPatterns: map[string]string{
				"email": "^.+$",
			},
		})
		require.NoError(t, err)
		userIDs := []string{users[0].ID.Hex(), users[2].ID.Hex()}
		require.Equal(t, len(userIDs), len(resp))
		for i := range resp {
			require.Contains(t, userIDs, resp[i].ID)
		}
		// check all with the provided email value.
		resp, err = st.Find(context.Background(), model.UserFindFilter{
			MetaPatterns: map[string]string{
				"email": "test2@gmail.com",
			},
		})
		require.NoError(t, err)
		userIDs = []string{users[2].ID.Hex()}
		require.Equal(t, len(userIDs), len(resp))
		for i := range resp {
			require.Contains(t, userIDs, resp[i].ID)
		}
	})
	t.Run("all ok (ids + status + meta in the filter)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := []MongoUser{
			{
				ID:     primitive.NewObjectID(),
				Status: "some status",
				Meta: map[string]interface{}{
					"email": "test@gmail.com",
				},
			},
			{
				ID:     primitive.NewObjectID(),
				Status: "some status 2",
				Meta: map[string]interface{}{
					"email": "test2@gmail.com",
				},
			},
			{
				ID:     primitive.NewObjectID(),
				Status: "some status 3",
				Meta: map[string]interface{}{
					"email": "test3gmail.com",
				},
			},
			{
				ID:     primitive.NewObjectID(),
				Status: "some status 3",
				Meta: map[string]interface{}{
					"email": "test4@gmail.com",
				},
			},
		}
		docs := make([]interface{}, len(users))
		for i := range users {
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			IDs:      []string{users[1].ID.Hex(), users[2].ID.Hex(), users[3].ID.Hex()},
			Statuses: []string{users[0].Status, users[2].Status, users[3].Status},
			MetaPatterns: map[string]string{
				"email": "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
			},
		})
		require.NoError(t, err)
		userIDs := []string{users[3].ID.Hex()}
		require.Equal(t, len(userIDs), len(resp))
		for i := range resp {
			require.Contains(t, userIDs, resp[i].ID)
		}
	})
	t.Run("all ok (with offset)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := make([]MongoUser, 20)
		docs := make([]interface{}, len(users))
		for i := range users {
			users[i] = MongoUser{
				ID: primitive.NewObjectID(),
			}
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		offset := int64(5)
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			Offset: &offset,
		})
		require.NoError(t, err)
		require.Equal(t, int64(len(users))-offset, int64(len(resp)))
		users = users[offset:]
		for i := range users {
			var found bool
			for j := range resp {
				if resp[j].ID == users[i].ID.Hex() {
					found = true
					break
				}
			}
			require.True(t, found)
		}
	})
	t.Run("all ok (with limit)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := make([]MongoUser, 10)
		docs := make([]interface{}, len(users))
		for i := range users {
			users[i] = MongoUser{
				ID: primitive.NewObjectID(),
			}
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		limit := int64(5)
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			Limit: &limit,
		})
		require.NoError(t, err)
		require.Equal(t, limit, int64(len(resp)))
		users = users[:limit]
		for i := range users {
			var found bool
			for j := range resp {
				if resp[j].ID == users[i].ID.Hex() {
					found = true
					break
				}
			}
			require.True(t, found)
		}
	})
	t.Run("all ok (with limit and offset)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		ctx := context.Background()
		users := make([]MongoUser, 20)
		docs := make([]interface{}, len(users))
		for i := range users {
			users[i] = MongoUser{
				ID: primitive.NewObjectID(),
			}
			docs[i] = users[i]
		}
		_, err := st.userCollection.InsertMany(ctx, docs)
		require.NoError(t, err)
		limit := int64(5)
		offset := int64(3)
		resp, err := st.Find(context.Background(), model.UserFindFilter{
			Limit:  &limit,
			Offset: &offset,
		})
		require.NoError(t, err)
		require.Equal(t, limit, int64(len(resp)))
		users = users[offset : limit+offset]
		for i := range users {
			var found bool
			for j := range resp {
				if resp[j].ID == users[i].ID.Hex() {
					found = true
					break
				}
			}
			require.True(t, found)
		}
	})
}

func createTestMongoStorage(t *testing.T) *MongoStorage {
	st, err := NewMongoStorage(context.Background(), testConnection, "test-db")
	require.NoError(t, err)
	require.NotNil(t, st)
	return st
}

func clearMongoStorage(t *testing.T, st *MongoStorage) {
	err := st.userCollection.Drop(context.Background())
	require.NoError(t, err)
}
