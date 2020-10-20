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

// NewMongoStorage returns new MongoStorage instance.
func NewMongoStorage(ctx context.Context, connString, dbName string) (*MongoStorage, error) {
	db, err := commonStorage.NewMongo(ctx, connString, dbName)
	if err != nil {
		return nil, err
	}

	userColl, err := db.Collection(ctx, userCollection, mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create %s collection", userColl)
	}

	return &MongoStorage{
		db:             db,
		userCollection: userColl,
	}, nil
}

// Disconnect breaks storage connection.
func (s *MongoStorage) Disconnect(ctx context.Context) error {
	return s.Disconnect(ctx)
}

// Add adds a new user.
func (s *MongoStorage) Add(ctx context.Context, user User) (*User, error) {
	mUser, err := user.ToMongoUser()
	if err != nil {
		return nil, err
	}

	res, err := s.userCollection.InsertOne(ctx, mUser)
	if err != nil {
		return nil, err
	}

	insertedID := res.InsertedID.(primitive.ObjectID)
	mUser.ID = &insertedID

	return mUser.ToUser(), nil
}
