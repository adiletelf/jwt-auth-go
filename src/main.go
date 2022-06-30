package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/handler"
	"github.com/adiletelf/jwt-auth-go/internal/middleware"
	"github.com/adiletelf/jwt-auth-go/internal/repository"
	"github.com/adiletelf/jwt-auth-go/internal/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	collection, err := util.GetCollection(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	cleanup(collection)

	tr := repository.NewTokenRepo(ctx, cfg, collection)
	h := handler.New(tr)
	r := gin.Default()

	configureRoutes(r, h)
	r.Run(cfg.ListenAddress)
}

func cleanup(collection *mongo.Collection) {
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

func configureRoutes(r *gin.Engine, h *handler.Handler) {
	r.GET("/generate", h.Generate)
	r.GET("/refresh", h.Refresh)

	// can't access without accessToken
	protected := r.Group("/api")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
}
