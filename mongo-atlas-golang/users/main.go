package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	path := request.Path
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if client == nil {
		return handleError(500, errors.New("no connection to DB"))
	}

	switch path {
	case "/users":
		log.Print("pathparameters: ", request.PathParameters)
		userID := request.PathParameters["id"]
		return getUserByID(userID)
	case "/health":
		return healthCheck(ctx)
	case "/ping":
		return ping()
	default:
		//return defaultReturn(request)
		userID := "5ce930b307a444000179a4e0"
		return getUserByID(userID)

	}

}

func getUserByID(userID string) (events.APIGatewayProxyResponse, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		handleError(400, errors.New("invalid id"))
	}

	var user User
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return handleError(404, errors.New("not found"))
		}
		return handleError(500, err)
	}

	return httpResponse(user)
}

func healthCheck(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	if err := client.Ping(ctx, nil); err != nil {
		return handleError(500, err)
	}
	return events.APIGatewayProxyResponse{
		Body:       string("ok"),
		StatusCode: 200,
	}, nil
}

func ping() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       string("pong"),
		StatusCode: 200,
	}, nil
}

func main() {
	if client == nil {
		connectToMongoDB()
	}
	lambda.Start(handler)

}

func handleError(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: status,
	}, err
}

func httpResponse(res User) (events.APIGatewayProxyResponse, error) {
	response, err := json.Marshal(res)
	if err != nil {
		return handleError(500, err)
	}
	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, nil
}
