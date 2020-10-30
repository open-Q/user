package storage

import (
	"bytes"
	"context"
	"testing"

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
		user := User{
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
		user := User{
			ID:     idHex,
			Status: "some status",
			Meta: map[string][]byte{
				"hello": []byte("world"),
				"key":   {1, 2, 3},
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
			require.True(t, bytes.Equal(v, value))
		}
	})
}

func TestMongoUser_ToUser(t *testing.T) {
	id := primitive.NewObjectID()
	idHex := id.Hex()
	userMongo := MongoUser{
		ID:     id,
		Status: "some status",
		Meta: map[string][]byte{
			"hello": []byte("world"),
			"key":   {1, 2, 3},
		},
	}
	res := userMongo.ToUser()
	require.NotNil(t, res)
	require.Equal(t, User{
		ID: idHex,
		Meta: map[string][]byte{
			"hello": []byte("world"),
			"key":   {1, 2, 3},
		},
		Status: userMongo.Status,
	}, *res)
}

func TestMongoStorage_Add(t *testing.T) {
	t.Run("convertation error", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := "invalid"
		_, err := st.Add(context.Background(), User{
			ID: id,
		})
		require.Error(t, err)
		require.EqualError(t, err, "encoding/hex: invalid byte: U+0069 'i'")
	})
	t.Run("insert error (duplicate entry id)", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		id := primitive.NewObjectID().Hex()
		_, err := st.Add(context.Background(), User{
			ID: id,
		})
		require.NoError(t, err)
		_, err = st.Add(context.Background(), User{
			ID: id,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate key error collection")
	})
	t.Run("all ok", func(t *testing.T) {
		st := createTestMongoStorage(t)
		defer clearMongoStorage(t, st)
		user, err := st.Add(context.Background(), User{
			Status: "some status",
			Meta: map[string][]byte{
				"hello": []byte("world"),
				"key":   {1, 2, 3},
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
