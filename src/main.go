package main

import (
	"log"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
	r.Run(cfg.ListenAddress)
}
