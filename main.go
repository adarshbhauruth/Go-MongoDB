package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Name string
}

func main() {

	db := Database{Name: "devdb"}

	http.HandleFunc("/", db.handler)

	// http.HandleFunc("/", httpHandle)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func (database Database) handler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	splitted_path := strings.Split(path, "/")

	num, _ := strconv.Atoi(splitted_path[2])
	num2, _ := strconv.Atoi(splitted_path[3])

	ctx, _ := context.WithTimeout(r.Context(), 10*time.Second)

	client, _ := createMongoClient("mongodb://db", ctx)

	collection := client.Database(database.Name).Collection("posts")

	sum := add(num, num2)

	date_current := time.Now()

	formatted_date := fmt.Sprintf("%d/%s/%d", date_current.Day(), date_current.Month().String(), date_current.Year())

	post := Post{Title: "addition", Number1: num, Number2: num2, Sum: sum, CreatedAt: formatted_date}

	insert(client, ctx, collection, &post)

	post_data, _ := postToJSON(&post)
	//post_data := postToJSON(w, &post)

	fmt.Fprintf(w, "%s", post_data)
}

func add(x int, y int) int {
	return x + y
}

func extract_path(r *http.Request) string {
	return r.URL.Path
}

func split_path(r *http.Request) []string {
	path := r.URL.Path
	splitted_path := strings.Split(path, "/")
	return splitted_path
}

func string_to_int(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

func createMongoClient(uri string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func insert(client *mongo.Client, ctx context.Context, collection *mongo.Collection, post *Post) error {
	// Insert Documents
	docs := []interface{}{
		bson.D{
			{Key: "title", Value: post.Title},
			{Key: "number1", Value: post.Number1},
			{Key: "number2", Value: post.Number2},
			{Key: "sum", Value: post.Sum},
			{Key: "created_at", Value: post.CreatedAt},
		},
	}

	_, insertErr := collection.InsertOne(ctx, docs[0])
	if insertErr != nil {
		return insertErr
	}
	return nil
}

func postToJSON(post *Post) ([]byte, error) {
	b, err := json.Marshal(post)
	//err := json.NewEncoder(w).Encode(&post)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// func postToJSON(w io.Writer, post *Post) error {
// 	//b, err := json.Marshal(post)
// 	err := json.NewEncoder(w).Encode(&post)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
