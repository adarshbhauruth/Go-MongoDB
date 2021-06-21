package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestPostServiceInsert(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := createMongoClient("mongodb://db", ctx)

	defer client.Disconnect(ctx)

	collection := client.Database("testdb").Collection("posts")
	collection.Drop(ctx)

	post := Post{Sum: 8}

	postService := PostService{collection: collection}

	err = postService.insert(&post, ctx)

	if err != nil {
		t.Fatalf("Expecting no error; got %v", err)
	}

	if collection == nil {
		t.Fatal("Could not find collection")
	}

	if err != nil {
		t.Fatalf("could not insert to db: %v", err)
	}

	// Iterate a cursor

	cur, curErr := collection.Find(ctx, bson.D{})

	if curErr != nil {
		panic(curErr)
	}

	defer cur.Close(ctx)

	var posts []Post
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}

	if len(posts) == 0 {
		t.Fatal("Expected post length to be 1; got 0")
	}

	actualPost := posts[0]

	if actualPost.Sum != post.Sum {
		t.Fatalf("Expected %d\nGot %d", post.Sum, actualPost.Sum)
	}

}

type Controller struct {
	postService PostServiceInterface
}

func (c Controller) handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := Post{Sum: 8}
	c.postService.insert(&post, ctx)
}

type MockPostService struct {
	counter int
}

func (p *MockPostService) insert(post *Post, ctx context.Context) error {
	p.counter++
	return nil
}

type PostServiceInterface interface {
	insert(post *Post, ctx context.Context) error
}

func TestHandler(t *testing.T) {

	mockPostService := MockPostService{}

	controller := Controller{postService: &mockPostService}
	r := httptest.NewRequest("GET", "/addition/3/5", nil)
	w := httptest.NewRecorder()
	controller.handler(w, r)
	if mockPostService.counter != 1 {
		t.Fatalf("Expected 1 got %d", mockPostService.counter)
	}
}
