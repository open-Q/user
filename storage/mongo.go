package storage

import (
	"context"

	commonStorage "github.com/open-Q/common/golang/storage"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollection = "user"
)

// MongoStorage represents mongo storage model.
type MongoStorage struct {
	db             *commonStorage.MongoStorage
	userCollection *commonStorage.MongoCollection
}

// MongoUser represents user mongo storage model.
type MongoUser struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Status string             `bson:"status"`
	Meta   map[string][]byte  `bson:"meta,omitempty"`
}

// NewMongoStorage returns new MongoStorage instance.
func NewMongoStorage(ctx context.Context, connString, dbName string) (*MongoStorage, error) {
	db, err := commonStorage.NewMongo(ctx, connString, dbName)
	if err != nil {
		return nil, err
	}

	userColl, err := db.Collection(ctx, userCollection, mongo.IndexModel{
		Keys: bson.M{
			"meta.0": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create %s collection", userCollection)
	}

	return &MongoStorage{
		db:             db,
		userCollection: userColl,
	}, nil
}

// Disconnect breaks storage connection.
func (s *MongoStorage) Disconnect(ctx context.Context) error {
	return s.db.Disconnect(ctx)
}

// Add adds a new user.
func (s *MongoStorage) Add(ctx context.Context, user User) (*User, error) {
	mUser, err := NewMongoUser(user)
	if err != nil {
		return nil, err
	}

	res, err := s.userCollection.InsertOne(ctx, mUser)
	if err != nil {
		return nil, err
	}

	mUser.ID = res.InsertedID.(primitive.ObjectID)

	return mUser.ToUser(), nil
}

// ToUser converts MongoUser model to User model.
func (m MongoUser) ToUser() *User {
	user := User{
		Status: m.Status,
		Meta:   m.Meta,
	}
	if !m.ID.IsZero() {
		user.ID = m.ID.Hex()
	}
	return &user
}

// NewMongoUser converts User model to MongoUser model.
func NewMongoUser(u User) (*MongoUser, error) {
	user := MongoUser{
		Status: u.Status,
		Meta:   u.Meta,
	}
	if u.ID != "" {
		id, err := primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return nil, err
		}
		user.ID = id
	}

	return &user, nil
}
