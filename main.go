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

	// establish connection to the server
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("An error occured with database connection")
		os.Exit(1)
	}

	// create context to be used by the connection
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	// make actual connection to the server, while also checking for errors
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Connection to db has failed")
	}

	// close the connection when this function exists
	defer client.Disconnect(ctx)

	// create database and collections
	quickstart := client.Database("quickstart")
	podcast := quickstart.Collection("podcast")
	episode := quickstart.Collection("episode")

	// create a new single document
	createdPodcast, err := podcast.InsertOne(ctx, bson.D{
		{Key: "title", Value: "The Joe Rogan Podcast"},
		{"author", "Joe Rogan"},
		{"medium", "Spotify"},
	})

	if err != nil {
		log.Fatal("Error while creating Podcast")
	}

	// insert multiple documents at a go
	createdEpisodes, err := episode.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "podcast", Value: createdPodcast.InsertedID},
			{Key: "title", Value: "Episode #1 - We Begin"},
			{Key: "description", Value: "This is the first episode in this podcast series"},
			{Key: "duration", Value: 25},
		},
		bson.D{
			{"podcast", createdPodcast.InsertedID},
			{"title", "The New Arena"},
			{"description", "This will help folks enter new arenas"},
			{"duration", 37},
			{"tags", bson.A{"though", "experiment", "era", "change"}},
		},
	})

	fmt.Println("Podcast ID: ", createdPodcast.InsertedID)
	fmt.Println("Episodes ID: ", createdEpisodes.InsertedIDs)
}
