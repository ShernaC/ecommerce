package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"users/config"
	"users/grpc/resolver"
	"users/middleware"
	"users/router"
	"utils/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	defaultPort     = "8080"
	defaultGRPCPort = "50051"
)

func init() {
	config.ConnectDB()
	config.SyncDB()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		port = defaultGRPCPort
	}

	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		r := gin.New()
		r.Use(
			gin.Recovery(),
			middleware.AuthMiddleware(),
			middleware.CORSMiddlewware(),
		)
		router.ApiRouter(r)

		log.Println("Listen and serve at http://localhost:" + port)
		r.Run(":" + port)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()

		user.RegisterUserServer(s, &resolver.Server{})
		fmt.Println("running at " + port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()
}
