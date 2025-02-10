package config

import (
	"github.com/darkphotonKN/journey-through-midnight/internal/server"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *server.Server) *gin.Engine {
	r := gin.Default()

	r.GET("/ws", server.HandleMatchConn)

	return r
}
