package main

import (
	"fmt"
	"log"
	"os"

	"github.com/darkphotonKN/journey-through-midnight/internal/config"
	"github.com/darkphotonKN/journey-through-midnight/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// env setup
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// db setup
	db := config.InitDB()
	defer db.Close()

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server := server.NewServer(port)

	// concurrently init messagehub listen-response loop
	go server.MessageHub()

	// routes setup
	routes := config.SetupRoutes(server)

	fmt.Printf("Server listening on port %s.\n", port)

	err := routes.Run(server.ListenAddr)

	if err != nil {
		log.Panic("Unable to start server. Err:", err)
	}
}
