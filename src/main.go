package main

import (
	"context"
	"log"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/handler"
	"github.com/adiletelf/jwt-auth-go/internal/repository"
	"github.com/adiletelf/jwt-auth-go/internal/util"
	"github.com/gin-gonic/gin"
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
	defer collection.Database().Drop(ctx)
	tr := repository.NewTokenRepo(ctx, cfg, collection)
	h := handler.New(tr)
	println(h)

	r := gin.Default()
	configureRoutes(r, h)
	r.Run(cfg.ListenAddress)
}

func configureRoutes(r *gin.Engine, h *handler.Handler) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	r.GET("/generate", h.Generate)
	r.GET("/refresh", h.Refresh)
}
