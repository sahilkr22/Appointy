package main

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// user represents data about a record user.
type user struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type posts struct {
	Id      string `json:"id"`
	Caption string `json:"caption"`
	Url     string `json:"url"`
}

func main() {
	router := gin.Default()
	router.GET("/users/:id", getuser)
	router.GET("/posts/users/:id", getuserposts)
	router.GET("/posts/:id", getposts)
	router.POST("/users", postuser)
	router.POST("/posts", postposts)

	router.Run("localhost:80")
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// getuser responds with the list of all users as JSON.
func getuser(c *gin.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://SAHIL:sah123il@appointy.giqnb.mongodb.net/Appointy?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	appointyDatabase := client.Database("Appointy")
	usersCollection := appointyDatabase.Collection("Users")

	var suser bson.M
	if err = usersCollection.FindOne(ctx, bson.M{"id": c.Params.ByName("id")}).Decode(&suser); err != nil {
		log.Fatal(err)
	}
	fmt.Println(suser)
	c.IndentedJSON(http.StatusOK, suser)
}

// postuser adds an user from JSON received in the request body.
func postuser(c *gin.Context) {
	var newuser user

	// Call BindJSON to bind the received JSON to
	// newuser.
	if err := c.BindJSON(&newuser); err != nil {
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://SAHIL:sah123il@appointy.giqnb.mongodb.net/Appointy?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	appointyDatabase := client.Database("Appointy")
	usersCollection := appointyDatabase.Collection("Users")

	usersResult, err := usersCollection.InsertOne(ctx, bson.D{
		{Key: "id", Value: newuser.Id},
		{Key: "name", Value: newuser.Name},
		{Key: "email", Value: newuser.Email},
		{Key: "password", Value: GetMD5Hash(newuser.Password)},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(usersResult.InsertedID)

	// Add the new user to the slice.

	c.IndentedJSON(http.StatusCreated, "User Added to the list")
}

func postposts(c *gin.Context) {
	var newpost posts

	// Call BindJSON to bind the received JSON to
	// newuser.
	if err := c.BindJSON(&newpost); err != nil {
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://SAHIL:sah123il@appointy.giqnb.mongodb.net/Appointy?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	appointyDatabase := client.Database("Appointy")
	postsCollection := appointyDatabase.Collection("Posts")

	postsResult, err := postsCollection.InsertOne(ctx, bson.D{
		{Key: "id", Value: newpost.Id},
		{Key: "caption", Value: newpost.Caption},
		{Key: "url", Value: newpost.Url},
		{Key: "date", Value: time.Now()},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(postsResult.InsertedID)

	// Add the new user to the slice.

	c.IndentedJSON(http.StatusCreated, newpost)
}

func getposts(c *gin.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://SAHIL:sah123il@appointy.giqnb.mongodb.net/Appointy?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	appointyDatabase := client.Database("Appointy")
	postsCollection := appointyDatabase.Collection("Posts")

	var spost bson.M
	objID, _ := primitive.ObjectIDFromHex(c.Params.ByName("id"))
	if err = postsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&spost); err != nil {
		log.Fatal(err)
	}
	fmt.Println(spost)
	c.IndentedJSON(http.StatusOK, spost)
}

func getuserposts(c *gin.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://SAHIL:sah123il@appointy.giqnb.mongodb.net/Appointy?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	appointyDatabase := client.Database("Appointy")
	postsCollection := appointyDatabase.Collection("Posts")

	cursor, err := postsCollection.Find(ctx, bson.M{"id": c.Params.ByName("id")})
	if err != nil {
		log.Fatal(err)
	}
	var spost []bson.M
	if err = cursor.All(ctx, &spost); err != nil {
		log.Fatal(err)
	}
	fmt.Println(spost)
	c.IndentedJSON(http.StatusOK, spost)
}
