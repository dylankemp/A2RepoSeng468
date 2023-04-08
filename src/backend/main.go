package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	//"os"
	"time"
	"context"

	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/handlers"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func connect() (*sql.DB, error) {
// 	bin, err := ioutil.ReadFile("/run/secrets/db-password")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sql.Open("mysql", fmt.Sprintf("root:%s@tcp(db:3306)/example", string(bin)))
// }

// func blogHandler(w http.ResponseWriter, r *http.Request) {
// 	db, err := connect()
// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("SELECT title FROM blog")
// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	var titles []string
// 	for rows.Next() {
// 		var title string
// 		err = rows.Scan(&title)
// 		titles = append(titles, title)
// 	}
// 	json.NewEncoder(w).Encode(titles)
// }

// func main2() {
// 	log.Print("Prepare db...")
// 	if err := prepare(); err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Print("Listening 8000")
// 	r := mux.NewRouter()
// 	r.HandleFunc("/", blogHandler)
// 	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
// }

// func prepare() error {
// 	db, err := connect()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	for i := 0; i < 60; i++ {
// 		if err := db.Ping(); err == nil {
// 			break
// 		}
// 		time.Sleep(time.Second)
// 	}

// 	if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
// 		return err
// 	}

// 	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id int NOT NULL AUTO_INCREMENT, title varchar(255), PRIMARY KEY (id))"); err != nil {
// 		return err
// 	}

// 	for i := 0; i < 5; i++ {
// 		if _, err := db.Exec("INSERT INTO blog (title) VALUES (?);", fmt.Sprintf("Blog post #%d", i)); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }


type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type Comment struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type Notification struct {
	PostNotification    string `json:"post"`
	CommentNotification string `json:"comment"`
}

type User struct {
	Name     string   `json:"name" bson:"name"`
	Username string   `json:"username" bson:"username"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	DoB      string   `json:"dob" bson:"dob"`
	Friends  []string `json:"friends" bson:"friends"`
}

var client *mongo.Client

func createUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	var newUser User

	json.NewDecoder(r.Body).Decode(&newUser)

	collection := client.Database("seng468").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, newUser)

	json.NewEncoder(w).Encode(result)
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	var users []User

	collection := client.Database("seng468").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":} "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":} "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(users)
}

func createPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	var newPost Post

	json.NewDecoder(r.Body).Decode(&newPost)

	collection := client.Database("seng468").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, newPost)

	// Redis setup
	conn, err := redis.Dial("tcp", "localhost:6379")
	checkError(err)
	defer conn.Close()
	_, err = conn.Do(
		"HMSET",
		newPost,
	)

	json.NewEncoder(w).Encode(result)
}

func getPosts(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	var posts []Post

	collection := client.Database("seng468").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":} "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":} "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func createComment(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	var newComment Comment

	json.NewDecoder(r.Body).Decode(&newComment)

	collection := client.Database("seng468").Collection("comments")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, newComment)

	json.NewEncoder(w).Encode(result)
}

func getComments(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")

	var comments []Comment

	collection := client.Database("seng468").Collection("comments")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":} "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment Comment
		cursor.Decode(&comment)
		comments = append(comments, comment)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":} "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(comments)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Starting the app")

	// Mongodb setup
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/comments", getComments).Methods("GET")
	router.HandleFunc("/comments", createComment).Methods("POST")
	http.ListenAndServe(":8089", router)
}
