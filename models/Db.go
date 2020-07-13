package models

import (
	"context"
	"fmt"
	"log"
	"os"
	_ "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var User *mongo.Collection
var ctx = context.TODO()

func init() {
	fmt.Println("Hiiii")
	envUri, ok := os.LookupEnv("MongoDB_URI")
	//mongoContext,_ := context.WithTimeout(context.Background(), 15 * time.Second)
	Uri := envUri
	if !ok{
		//log.Print("unable to connect to remote Database")
		Uri = "mongodb://localhost:27017"
	}
	clientOptions := options.Client().ApplyURI(Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	//Blog := client.Database("Blog-Api").Collection("blog")
	User = client.Database("Blog-Api").Collection("user")
}
