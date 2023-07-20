package main

import (
	api_v1 "gpt_presets_backend/internal/api/v1"
	"gpt_presets_backend/internal/database"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.Connect()
	defer database.Disconnect(db)

	r := gin.Default()

	api_v1.Routes(r, db)

	r.Run()
}
