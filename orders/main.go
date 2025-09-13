package main

import (
	"log"
	"orders/config"
	"orders/middleware"
	"orders/router"
	"os"

	"github.com/gin-gonic/gin"
)

var defaultPort = "8080"

func init() {
	config.ConnectDB()
	config.SyncDB()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := gin.New()
	r.Use(
		gin.Recovery(),
		middleware.AuthMiddleware(),
		middleware.CORSMiddlewware(),
	)
	router.ApiRouter(r)

	log.Println("Listen and serve at http://localhost:" + port)
	r.Run(":" + port)
}
