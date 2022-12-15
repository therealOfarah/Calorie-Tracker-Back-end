package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
func main()  {
	port := os.Getenv("PORT")
	if port ==""{
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
}