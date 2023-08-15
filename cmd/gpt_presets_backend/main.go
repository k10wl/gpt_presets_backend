package main

import (
	"log"

	api_v1 "gpt_presets_backend/internal/app/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gpt_presets_backend/internal/app/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	db := database.Connect()
	defer database.Disconnect(db)

	api_v1.Routes(r, db)

	r.Run()
}
