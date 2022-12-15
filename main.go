package main

import (
	"os"
	"github.com/therealofarah/go-calorie-tracker/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)
func main()  {
	port := os.Getenv("PORT")
	if port ==""{
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	//routes
	router.POST("/entry/create", routes.AddEntry)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entry/:id",routes.GetEntryById)
	router.GET("/ingredient/:ingredient", routes.GetEntriesByIngredient)
	router.PUT("/entry/update/:id",routes.UpdateEntry)
	router.PUT("/ingredient/update/:id",routes.UpdateIngredient)
	router.DELETE("/entry/delete/:id",routes.DeleteEntry)
	router.Run(":"+port)
}