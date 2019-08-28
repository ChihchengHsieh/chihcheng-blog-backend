package models

import (
	"DiscussionBoard/utils"
	"blog/databases"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - Schema for User
type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

var projectionForRemovingPassword = bson.D{
	{"password", 0},
}

// AddUser - Add a user to the database
func AddUser(inputUser *User) (interface{}, error) {
	result, err := databases.DB.Collection("user").InsertOne(context.TODO(), inputUser)
	return result.InsertedID, err
}

// UpdateUserByID - update the user by given fields
func UpdateUserByID(id string, updateDetail map[string]interface{}) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := databases.DB.Collection("user").UpdateOne(context.TODO(), bson.M{"id": oid}, bson.M{"$set": updateDetail})

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

// FindUserByID - Find the user through the given ID
func FindUserByID(id string) (*User, error) {
	var user *User
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	err = databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"_id": oid}, options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByEmail - Find the user through given Email
func FindUserByEmail(email string) (*User, error) {
	var user *User
	err := databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email}, options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CheckingTheAuth - if the user is returned, then the auth is valid
func CheckingTheAuth(email string, password string) (*User, error) {
	var user User
	err := databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUsers - this function will get multiple users without the password field
func FindUsers(filterDetail bson.M) ([]*User, error) {
	var users []*User
	result, err := databases.DB.Collection("user").Find(context.TODO(), filterDetail,
		options.Find().SetProjection(projectionForRemovingPassword))

	if err != nil {
		return nil, err
	}
	defer result.Close(context.TODO())

	for result.Next(context.TODO()) {
		var elem User
		err := result.Decode(&elem)
		utils.ErrorChecking(err)
		users = append(users, &elem)
	}
	return users, nil
}
