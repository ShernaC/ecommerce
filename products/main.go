package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"products/config"
	"products/grpc/resolver"
	"products/router"
	"sync"
	"utils/middleware"
	"utils/product"

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

	grpcPort := os.Getenv("PRODUCT_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = defaultGRPCPort
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

		product.RegisterProductServer(s, resolver.Server{})
		fmt.Println("running at " + grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()
}
