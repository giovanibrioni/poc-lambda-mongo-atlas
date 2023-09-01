package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

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
	col    = "users"
)

type User struct {
	Base
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	First   string             `bson:"first" json:"first" binding:"required"`
	Last    string             `bson:"last" json:"last" binding:"required"`
	Email   string             `bson:"email" json:"email" binding:"required"`
	Status  string             `bson:"status" json:"status" binding:"required"`
	City    string             `bson:"city" json:"city"`
	Country string             `bson:"country" json:"country"`
	Age     int                `bson:"age" json:"age" binding:"required"`
}

type Base struct {
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}

func connectToDB(ctx context.Context) {
	if client != nil {
		log.Printf("client already connected")
		return
	}
	var connectionError error
	opts := options.Client()
	secondary := readpref.Secondary()
	opts.SetReadPreference(secondary)
	opts.SetServerSelectionTimeout(5 * time.Second)
	opts.SetConnectTimeout(2 * time.Second)
	opts.ApplyURI(os.Getenv("MONGODB_URI"))

	client, connectionError = mongo.Connect(ctx, opts)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	log.Printf("connected to the database.")
	collection = client.Database(dbName).Collection(col)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/v1/user", listUsers)
	router.GET("/v1/user/:id", getUserByID)
	router.POST("/v1/user", createUser)
	router.PUT("/v1/user/:id", updateUser)
	router.DELETE("/v1/user/:id", deleteUser)
	router.GET("/health", healthCheck)
	router.GET("/ping", ping)

	return router

}

func listUsers(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Calculate the offset based on the page number and limit
	offset := (page - 1) * limit

	userList := make([]*User, 0)
	filter := bson.D{}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, opts)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(context.TODO(), &userList); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       userList,
		"page":       page,
		"limit":      limit,
		"totalCount": totalCount,
	})
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
		log.Printf("Error finding user %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = primitive.NewObjectID()
	user.Status = "active"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if _, err := collection.InsertOne(context.Background(), user); err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func updateUser(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UpdatedAt = time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.M{"_id": id}
	update := bson.D{{Key: "$set", Value: user}}
	err = collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func deleteUser(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	res, err := collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if res.DeletedCount < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func healthCheck(c *gin.Context) {
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Print(err)
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
		log.Printf("connecting to the database...")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		connectToDB(ctx)
		defer cancel()
	}
	lambdaName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if lambdaName != "" {
		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {
		port := "8080"
		log.Println("start server and listen on port:", port)
		log.Fatal(http.ListenAndServe(":"+port, setupRouter()))
	}

}
