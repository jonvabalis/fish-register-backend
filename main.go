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

	r.POST("/register", fishApi.Register)
	r.PATCH("/change-login", fishApi.ChangeLogin)
	r.GET("/users", fishApi.GetUsers)
	r.DELETE("/users", fishApi.DeleteUser)

	r.GET("/locations", fishApi.GetLocations)
	r.POST("/locations", fishApi.InsertLocation)
	r.PATCH("/locations", fishApi.PatchLocation)

	r.GET("/species", fishApi.GetAllSpecies)
	r.POST("/species", fishApi.InsertSpecies)
	r.PATCH("/species", fishApi.PatchSpecies)

	r.GET("/locations-species", fishApi.GetAllSpeciesByLocation)
	r.POST("/locations-species", fishApi.InsertSpeciesToLocation)
	r.DELETE("/locations-species", fishApi.DeleteSpeciesFromLocation)

	err = r.Run(":1111")
	if err != nil {
		return
	}
}
