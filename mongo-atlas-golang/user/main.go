package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/apex/gateway/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
)

const (
	dbName = "test"
	col    = "usuario"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"nome" json:"nome"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"senha" json:"-"`
	Status    string             `bson:"status" json:"status"`
	Role      []string           `bson:"perfil" json:"role"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

func connectToMongoDB() {
	if client != nil {
		return
	}
	var connectionError error
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	opts := options.Client()
	opts.SetServerSelectionTimeout(1 * time.Second)
	opts.SetConnectTimeout(2 * time.Second)
	opts.SetSocketTimeout(2 * time.Second)
	opts.ApplyURI(os.Getenv("MONGODB_URI"))

	client, connectionError = mongo.Connect(ctx, opts)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	collection = client.Database(dbName).Collection(col)
	defer cancel()

}

func setupRouter() *gin.Engine {
	log.Printf("Gin cold start")
	router := gin.Default()

	router.GET("/users/:id", getUserByID)
	router.GET("/health", healthCheck)
	router.GET("/ping", ping)

	return router

}

func getUserByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user User
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func healthCheck(c *gin.Context) {
	if err := client.Ping(context.TODO(), nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	if client == nil {
		connectToMongoDB()
	}
	lambdaName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if lambdaName != "" {
		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {
		port := "8080"
		fmt.Println("start server and listen on port: ", port)
		log.Fatal(http.ListenAndServe(":"+port, setupRouter()))
	}

}
