package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(ctx context.Context, cfg *config.Config) (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(cfg.DB.ConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(cfg.DB.Name).Collection(cfg.DB.CollectionName)
	return collection, nil
}

func Cleanup(collection *mongo.Collection) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Performing cleanup...")
		err := collection.Database().Drop(context.TODO())
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
}
