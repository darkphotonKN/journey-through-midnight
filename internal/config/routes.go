package config

import (
	"fmt"

	"github.com/darkphotonKN/journey-through-midnight/internal/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *server.Server) *gin.Engine {
	router := gin.Default()

	// NOTE: debugging middleware
	router.Use(func(c *gin.Context) {
		fmt.Println("Incoming request to:", c.Request.Method, c.Request.URL.Path, "from", c.Request.Host)
		c.Next()
	})

	// TODO: CORS for development, remove in PROD
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization",
			"Accept",
			"X-Requested-With",
			"X-CSRF-Token",
			"Cache-Control",
			"*", // TODO: remove in prod
		},
		AllowCredentials: true,
	}))

	router.GET("/ws", server.HandleMatchConn)

	api := router.Group("/api")
	api.POST("/signup", server.HandleMatchConn)

	return router
}
