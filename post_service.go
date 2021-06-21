package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Post struct {
	Title   string `bson:"title,omitempty"`
	Number1 int    `bson:"number1,omitempty"`
	Number2 int    `bson:"number2,omitempty"`
	Sum     int    `bson:"sum,omitempty"`
	//CreatedAt time.Time `bson:"created_at"`
	CreatedAt string `bson:"created_at"`
}

type PostService struct {
	collection *mongo.Collection
}

func (p PostService) insert(post *Post, ctx context.Context) error {

	_, err := p.collection.InsertOne(ctx, post)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return err
	}
	return nil
}
