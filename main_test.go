package main

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert(t *testing.T) {

	title := "addition"

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := createMongoClient("mongodb://db", ctx)

	collection := client.Database("testdb").Collection("posts")

	if collection == nil {
		t.Fatal("Could not find collection")
	}

	defer client.Disconnect(ctx)
	collection.Drop(ctx)

	post := Post{Title: title, Number1: 5, Number2: 7, Sum: 12, CreatedAt: "16/06/21"}

	err = insert(client, ctx, collection, &post)

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

	if posts[0].Title != title {
		t.Fatal("Title does not match input")
	}
}

func TestPostToJSON(t *testing.T) {
	expectedValue := `{"Title":"Addition","Number1":5,"Number2":7,"Sum":12,"CreatedAt":"16/June/2021"}`
	post := Post{Title: "Addition", Number1: 5, Number2: 7, Sum: 12, CreatedAt: "16/June/2021"}

	//err := postToJSON(os.Stdout, &post)
	b, err := postToJSON(&post)

	if err != nil {
		t.Fatal("Could not marshal POST")
	}
	if expectedValue != string(b) {
		t.Fatalf("Expected: %s\n Got: %s", expectedValue, string(b))
	}

}

// func TestGet(t *testing.T) {
// 	want := "Success!"
// 	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(200)
// 		w.Write([]byte(want))
// 	}))
// 	defer srv.Close()

// 	sut := &Curl{
// 		Client: srv.Client(),
// 		URL:    srv.URL,
// 	}

// 	got, err := sut.Get(map[string]string{})
// 	if err != nil {
// 		t.Errorf("Unexpected error on request: %s", err)
// 	}
// 	if got != want {
// 		t.Errorf("want %s, got %s", want, got)
// 	}
// }
