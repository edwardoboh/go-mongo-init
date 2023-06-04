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
		log.Fatal(err)
	}

	// create context to be used by the connection
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	// make actual connection to the server, while also checking for errors
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// close the connection when this function exists
	defer client.Disconnect(ctx)

	// create database and collections
	quickstart := client.Database("quickstart")
	// podcast := quickstart.Collection("podcast")
	episode := quickstart.Collection("episode")

	/*

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

	*/

	// Fetching all created records from a collection
	// cursor, err := episode.Find(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cursor.Close(ctx)

	// METHOD #1
	/*
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			log.Fatal(err)
		}

		for _, epi := range results {
			fmt.Println(epi)
			fmt.Println()
		}
	*/

	// Method #2
	/*
		for cursor.Next(ctx) {
			var epi bson.M
			if err2 := cursor.Decode(&epi); err2 != nil {
				log.Fatal(err2)
			}
			fmt.Println(epi)
		}
	*/

	// Fetching the first record from the collection
	// var firstEpisode bson.M
	// if err := episode.FindOne(ctx, bson.M{}).Decode(&firstEpisode); err != nil {
	// 	log.Fatal()
	// }

	// fmt.Println(firstEpisode)

	// NOTE - user bson.D when orders matters and bson.M when id doesn't

	// Filtering find operation to get specific document(s)
	// cursor, err := episode.Find(ctx, bson.M{"duration": bson.M{"$gte": 20, "$lte": 30}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var found []bson.M
	// if err := cursor.All(ctx, &found); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(found)

	// Ordering documents using Find options
	opts := options.Find()
	opts.SetSort(bson.D{{"duration", -1}, {"_id", 1}})
	sortedCursor, err := episode.Find(ctx, bson.D{
		{"duration", bson.D{
			{"$gte", 20},
			{"$lte", 50},
		}},
	}, opts)

	var sortedResp []bson.M

	if err := sortedCursor.All(ctx, &sortedResp); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sortedResp)
}
