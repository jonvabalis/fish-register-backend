package main

import (
	"log"

	"fish-register-backend/internal/db"
	"fish-register-backend/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("could not initialize db: %s", err)
	}
	defer dbConn.Close()

	log.Println("application started successfully")

	r := gin.Default()
	r.Use(cors.Default())

	fishApi := handlers.NewFishApi(dbConn)

	r.GET("/register", fishApi.Register)

	r.GET("/locations", fishApi.GetLocations)
	r.POST("/locations", fishApi.InsertLocation)
	r.PATCH("/locations", fishApi.PatchLocation)

	err = r.Run(":1111")
	if err != nil {
		return
	}
}
