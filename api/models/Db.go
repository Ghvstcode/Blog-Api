package models

import (
	"context"
	"os"
	_ "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	l "github.com/GhvstCode/Blog-Api/api/utils/logger"
)

var User *mongo.Collection
var Blog *mongo.Collection
var ctx = context.TODO()

func init() {
	l.InfoLogger.Println("Connecting to DB...")
	envUri, ok := os.LookupEnv("MongoDB_URI")

	Uri := envUri
	if !ok {
		l.WarningLogger.Println("Unable to connect to load connection URI from env file,connecting to local db!")
		//Uri = "mongodb://localhost:27017"
		Uri = "mongodb://mongo:27017"
	}
	clientOptions := options.Client().ApplyURI(Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.ErrorLogger.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		l.ErrorLogger.Fatal(err)
	}

	Blog = client.Database("Blog-Api").Collection("blog")
	User = client.Database("Blog-Api").Collection("user")
}
