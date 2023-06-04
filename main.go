package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadEnv() {
	godotenv.Load()
}

func main() {
	loadEnv()

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("An error occured with database connection")
		os.Exit(1)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Connection to db has failed")
	}

	defer client.Disconnect(ctx)

	names, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(names)
}
