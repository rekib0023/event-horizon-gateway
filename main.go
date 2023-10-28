// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	RegisterAuthEndpoints(r, gRpc)
	RegisterProfileEndpoints(r, gRpc)

	port := os.Getenv("AUTH_SVC")
	if port == "" {
		log.Fatal("AUTH_SVC environment variable not set")
	}

	log.Println("Starting server on :" + port + "...")
	r.Run(":" + port)
}
