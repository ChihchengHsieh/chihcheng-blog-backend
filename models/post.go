package models

import (
	"blog/databases"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post - The Logging Post
type Post struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Author      primitive.ObjectID `json:"author" bson:"author"`
	Title       string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"` // This should be a markdown
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	PrivatePost bool               `json:"privatePost" bson:"privatePost"`
}

// AddPost - Add a Post
func AddPost(inputPost Post) (interface{}, error) {
	result, err := databases.DB.Collection("post").InsertOne(context.TODO(), inputPost)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// UpdatePostByID - Update the post through given ID
func UpdatePostByID(id string, updateDetail map[string]interface{}) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	result, err := databases.DB.Collection("post").UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

// DeletePostByID - Delete the post through given ID
func DeletePostByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = databases.DB.Collection("post").DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return nil
}

// FindOnePost - by the given detail
func FindOnePost(filterDetail bson.M) (*Post, error) {
	var post *Post

	err := databases.DB.Collection("post").FindOne(context.TODO(), filterDetail).Decode(&post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

// FindPostByID - Find the post with given ID
func FindPostByID(id string) (*Post, error) {
	var post *Post

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = databases.DB.Collection("post").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

// FindPosts - Find all the posts with the given condition
func FindPosts(filterDetail bson.M, findOptions *options.FindOptions) ([]*Post, error) {
	var posts []*Post

	fmt.Println(filterDetail)

	result, err := databases.DB.Collection("post").Find(context.TODO(), filterDetail, findOptions)

	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Post
		err := result.Decode(&elem)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &elem)
	}

	return posts, nil
}

// FindPostsPagination - receiveing the pagination parameter to generate results
func FindPostsPagination(filterDetail bson.M, sortDetail bson.M, skip int, limit int) ([]*Post, error) {
	var posts []*Post

	fmt.Printf("filterDetail is %+v\n", filterDetail)

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(sortDetail)

	result, err := databases.DB.Collection("post").Find(context.TODO(),
		filterDetail,
		findOptions,
	)

	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Post
		err := result.Decode(&elem)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &elem)
	}

	return posts, nil

}
