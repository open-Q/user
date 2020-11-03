package storage

import (
	"context"
	"log"

	commonErrors "github.com/open-Q/common/golang/errors"
	commonStorage "github.com/open-Q/common/golang/storage"
	"github.com/open-Q/user/storage/model"
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
	ID     primitive.ObjectID     `bson:"_id,omitempty"`
	Status string                 `bson:"status"`
	Meta   map[string]interface{} `bson:"meta,omitempty"`
}

// NewMongoStorage returns new MongoStorage instance.
func NewMongoStorage(ctx context.Context, connString, dbName string) (*MongoStorage, error) {
	db, err := commonStorage.NewMongo(ctx, connString, dbName)
	if err != nil {
		return nil, err
	}

	userColl, err := db.Collection(ctx, userCollection)
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
func (s *MongoStorage) Add(ctx context.Context, user model.User) (*model.User, error) {
	mUser, err := NewMongoUser(user)
	if err != nil {
		return nil, commonErrors.NewStorageConvertError(err.Error())
	}

	res, err := s.userCollection.InsertOne(ctx, mUser)
	if err != nil {
		return nil, commonErrors.NewStorageInsertError(err.Error())
	}

	mUser.ID = res.InsertedID.(primitive.ObjectID)

	return mUser.ToUser(), nil
}

// Delete removes an existing user by ID.
func (s *MongoStorage) Delete(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return commonErrors.NewStorageConvertError(err.Error())
	}

	filter := bson.M{
		"_id": id,
	}

	res, err := s.userCollection.DeleteOne(ctx, filter)
	if res.DeletedCount == 0 && err == nil {
		err = errors.New("user not found")
	}
	if err != nil {
		return commonErrors.NewStorageDeleteError(err.Error())
	}

	return nil
}

// Update updates an existing user.
func (s *MongoStorage) Update(ctx context.Context, user model.User) (*model.User, error) {
	mUser, err := NewMongoUser(user)
	if err != nil {
		return nil, commonErrors.NewStorageConvertError(err.Error())
	}

	filter := bson.M{
		"_id": mUser.ID,
	}
	res, err := s.userCollection.ReplaceOne(ctx, filter, mUser)
	if res.MatchedCount == 0 && err == nil {
		err = errors.New("user not found")
	}
	if err != nil {
		return nil, commonErrors.NewStorageUpdateError(err.Error())
	}

	return mUser.ToUser(), nil
}

// Find finds users by filter.
func (s *MongoStorage) Find(ctx context.Context, filter model.UserFindFilter) ([]model.User, error) {
	mongoFilter, findOptions, err := createUserFindFilter(filter)
	if err != nil {
		return nil, err
	}

	cursor, err := s.userCollection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, commonErrors.NewStorageFindError(err.Error())
	}
	defer closeCursor(ctx, cursor)

	var users []MongoUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, commonErrors.NewStorageConvertError(err.Error())
	}

	foundUsers := make([]model.User, len(users))
	for i := range users {
		foundUsers[i] = *users[i].ToUser()
	}

	return foundUsers, nil
}

// ToUser converts MongoUser model to User model.
func (m MongoUser) ToUser() *model.User {
	user := model.User{
		Status: m.Status,
		Meta:   convertMeta(m.Meta),
	}
	if !m.ID.IsZero() {
		user.ID = m.ID.Hex()
	}
	return &user
}

// NewMongoUser converts User model to MongoUser model.
func NewMongoUser(u model.User) (*MongoUser, error) {
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

func convertMeta(meta map[string]interface{}) map[string]interface{} {
	for k, v := range meta {
		meta[k] = spreadPrimitives(v)
	}

	return meta
}

func spreadPrimitives(value interface{}) interface{} {
	switch value.(type) {
	case primitive.A:
		return spreadPrimitives([]interface{}(value.(primitive.A)))
	case map[string]interface{}:
		m := value.(map[string]interface{})
		for k, v := range m {
			m[k] = spreadPrimitives(v)
		}
		return m
	case []interface{}:
		aValues, ok := value.([]interface{})
		if !ok {
			return value
		}
		spreads := make([]interface{}, 0, len(aValues))
		for i := range aValues {
			spreads = append(spreads, spreadPrimitives(aValues[i]))
		}
		return spreads
	}

	return value
}

func createUserFindFilter(filter model.UserFindFilter) (bson.M, *options.FindOptions, error) {
	mongoFilter := bson.M{}

	if len(filter.IDs) != 0 {
		ids := make([]primitive.ObjectID, len(filter.IDs))
		for i := range filter.IDs {
			id, err := primitive.ObjectIDFromHex(filter.IDs[i])
			if err != nil {
				return nil, nil, commonErrors.NewStorageConvertError(err.Error())
			}
			ids[i] = id
		}
		mongoFilter["_id"] = bson.M{
			"$in": ids,
		}
	}

	if len(filter.Statuses) != 0 {
		mongoFilter["status"] = bson.M{
			"$in": filter.Statuses,
		}
	}

	if len(filter.MetaPatterns) != 0 {
		for k, v := range filter.MetaPatterns {
			mongoFilter["meta."+k] = bson.M{
				"$regex": primitive.Regex{
					Options: "i",
					Pattern: v,
				},
			}
		}
	}

	opts := options.Find()
	if filter.Offset != nil || filter.Limit != nil {
		opts.SetSort(bson.M{
			"_id": 1,
		})
	}
	if filter.Offset != nil {
		opts.SetSkip(int64(*filter.Offset))
	}
	if filter.Limit != nil {
		opts.SetLimit(int64(*filter.Limit))
	}

	return mongoFilter, opts, nil
}

func closeCursor(ctx context.Context, cursor *mongo.Cursor) {
	if err := cursor.Close(ctx); err != nil {
		log.Printf("could not close cursor: %v", err)
	}
}
