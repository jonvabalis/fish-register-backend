package main

import (
	"fish-register-backend/internal/db"
	"fish-register-backend/internal/handlers"
	"log"

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
	fishApi := handlers.NewFishApi(dbConn)

	r.GET("/register", fishApi.Register)

	err = r.Run(":1111")
	if err != nil {
		return
	}
}
