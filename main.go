package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/javiercanadillas/mongogo/mongosecrets"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection
var validate = validator.New()

type Book struct {
	Id              primitive.ObjectID `json:"id,omitempty"`
	Title           string             `json:"title,omitempty" validate:"required"`
	PageCount       int                `json:"pageCount,omitempty"`
	LongDescription string             `json:"longDescription,omitempty"`
	ISBN            string             `json:"isbn,omitempty"`
	Authors         []string           `json:"authors,omitempty"`
	Categories      []string           `json:"categories,omitempty"`
	Status          string             `json:"status,omitempty"`
}

type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// Create a new book
func createBook(c *gin.Context) {
	var book Book

	// Check request body
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Error in request",
			Data: map[string]interface{}{
				"data": err.Error(),
			}})
	}

	// Validate required fields
	if validationErr := validate.Struct(&book); validationErr != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: map[string]interface{}{
				"data": validationErr.Error(),
			},
		})
		return
	}

	newBook := Book{
		Id:              primitive.NewObjectID(),
		PageCount:       book.PageCount,
		Title:           book.Title,
		ISBN:            book.ISBN,
		Authors:         book.Authors,
		Categories:      book.Categories,
		Status:          book.Status,
		LongDescription: book.LongDescription,
	}

	result, err := coll.InsertOne(context.TODO(), newBook)

	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error inserting book",
			Data: map[string]interface{}{
				"data": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: map[string]interface{}{
			"data": result,
		},
	})
}

// Old reading function, to be deprecated by readBook, left here for
// backwards compatibility
func read(c *gin.Context) {
	title := c.Query("title")
	filter := bson.M{
		"title": title,
	}
	var result bson.M

	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		log.Fatalf("No book was found with the title %s\n", title)
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	c.JSON(http.StatusOK, result)
}

// Get book by title
func readBook(c *gin.Context) {
	title := c.Param("title")
	filter := bson.M{
		"title": title,
	}
	var book Book

	err := coll.FindOne(context.TODO(), filter).Decode(&book)

	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error searching for book %s", title),
			Data: map[string]interface{}{
				"data": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Got result for title %s", title),
		Data: map[string]interface{}{
			"data": book,
		},
	})
}

// Delete a book by title
func deleteBook(c *gin.Context) {
	title := c.Param("title")
	filter := bson.M{
		"title": title,
	}

	result, err := coll.DeleteOne(context.TODO(), filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error when trying to delete book %s", title),
			Data: map[string]interface{}{
				"data": err.Error(),
			},
		})
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound, UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data: map[string]interface{}{
				"data": fmt.Sprintf("Could not find a book titled %s", title),
			},
		})
	}

	c.JSON(http.StatusOK, UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"data": fmt.Sprintf("Success deleting book %s", title),
		},
	})
}

func main() {
	// Init Mongo connection with Connection URL
	log.Println("Trying to get connection URL from Cloud Secret Manager")
	uri := fmt.Sprintf("%s", mongosecrets.GetSecret("MongoConnURL", "latest"))
	if uri == "" {
		log.Println("Couldn't retrieve URL from Cloud Secret Manager. Trying local env vars.")
		uri = os.Getenv("MONGODB_URI")
		if uri == "" {
			log.Fatal("You must set your 'MONGODB_URI' env variable.")
		}
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	coll = client.Database("library").Collection("books")

	// Set gin-gonic in prod mode
	gin.SetMode(gin.ReleaseMode)	
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": "Hello from my basic REST API",
		})
	})
	r.GET("/book/:title", readBook)
	r.GET("/read", read)
	r.POST("/book", createBook)
	r.DELETE("/book/:title", deleteBook)
	r.Run("localhost:8080")
}

