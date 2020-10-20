package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents user storage model.
type User struct {
	ID     *string
	Email  string
	Status string
}

// MongoUser represents user mongo storage model.
type MongoUser struct {
	ID     *primitive.ObjectID `bson:"_id,omitempty"`
	Email  string              `bson:"email"`
	Status string              `bson:"status"`
}

// ToMongoUser converts User model to MongoUser model.
func (u User) ToMongoUser() (*MongoUser, error) {
	user := MongoUser{
		Email:  u.Email,
		Status: u.Status,
	}
	if user.ID != nil {
		id, err := primitive.ObjectIDFromHex(*u.ID)
		if err != nil {
			return nil, err
		}
		user.ID = &id
	}

	return &user, nil
}

// ToUser converts MongoUser model to User model.
func (m MongoUser) ToUser() *User {
	user := User{
		Email:  m.Email,
		Status: m.Status,
	}
	if m.ID != nil {
		id := m.ID.Hex()
		user.ID = &id
	}
	return &user
}
