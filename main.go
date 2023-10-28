// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rekib0023/event-horizon-gateway/api"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conn, err := grpc.Dial(os.Getenv("AUTH_SVC"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	r := gin.Default()

	gRpc := pb.NewAuthServiceClient(conn)

	authGroup := r.Group("/api/auth")
	api.RegisterAuthEndpoints(authGroup, gRpc)

	usersGroup := r.Group("/api/users")
	api.RegisterProfileEndpoints(usersGroup, gRpc)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	log.Println("Starting server on :" + port + "...")
	r.Run(":" + port)
}
