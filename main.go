package main

import (
	"log"
	"strings"
	"time"

	"fish-register-backend/internal/db"
	"fish-register-backend/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	fishApi := handlers.NewFishApi(dbConn)

	r.POST("/register", fishApi.Register)
	r.POST("/login", fishApi.Login)
	r.PATCH("/change-login", fishApi.ChangeLogin)

	// Routes requiring authorization
	pr := r.Group("/", withAuthorization())

	pr.GET("/users", fishApi.GetUsers)
	pr.DELETE("/users", fishApi.DeleteUser)

	pr.GET("/locations", fishApi.GetLocations)
	pr.POST("/locations", fishApi.InsertLocation)
	pr.PATCH("/locations", fishApi.PatchLocation)
	pr.DELETE("/locations", fishApi.DeleteLocation)

	pr.GET("/species", fishApi.GetAllSpecies)
	pr.POST("/species", fishApi.InsertSpecies)
	pr.PATCH("/species", fishApi.PatchSpecies)
	pr.DELETE("/species", fishApi.DeleteSpecies)

	pr.GET("/locations/:locationUUID/species", fishApi.GetAllSpeciesByLocation)
	pr.POST("/locations/species", fishApi.InsertSpeciesToLocation)
	pr.DELETE("/locations/species", fishApi.DeleteSpeciesFromLocation)

	pr.GET("/users/:userUUID/rods", fishApi.GetUserRods)
	pr.POST("/rods", fishApi.InsertRod)
	pr.PATCH("/rods", fishApi.PatchRod)
	pr.DELETE("/rods", fishApi.DeleteRod)

	pr.POST("/catches", fishApi.CreateCatch)
	pr.GET("/users/:userUUID/catches", fishApi.GetUserCatches)
	pr.PATCH("/catches", fishApi.UpdateUserCatch)
	pr.DELETE("/catches", fishApi.DeleteCatch)

	err = r.Run(":1111")
	if err != nil {
		return
	}
}

func withAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
			return []byte("test"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(int64); ok {
				if time.Now().Unix() > exp {
					c.JSON(401, gin.H{"error": "token expired"})
					c.Abort()
					return
				}
			}

			c.Next()
		}

		c.Abort()
	}
}
