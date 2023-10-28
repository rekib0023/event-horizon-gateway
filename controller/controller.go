package controller

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
	"google.golang.org/grpc"
)

type Controller struct {
	r    *gin.RouterGroup
	gRpc pb.AuthServiceClient
}

var controller *Controller

func Init() {
	controller.InitAuthController()
	controller.InitProfileController()
}

var e *gin.Engine

func initGin() {

}

func Start() {
	e = gin.Default()

	conn, err := grpc.Dial(os.Getenv("AUTH_SVC"), grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
	} else {
		defer conn.Close()
		gRpc := pb.NewAuthServiceClient(conn)
		apiGroup := e.Group("/api")
		controller = &Controller{
			r:    apiGroup,
			gRpc: gRpc,
		}
	}

	Init()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	log.Println("Starting server on :" + port + "...")
	e.Run(":" + port)
}

func POST(pattern string, handler gin.HandlerFunc) {
	controller.r.POST(pattern, handler)
}

func GET(pattern string, handler gin.HandlerFunc) {
	controller.r.GET(pattern, handler)
}

func PUT(pattern string, handler gin.HandlerFunc) {
	controller.r.PUT(pattern, handler)
}

func DELETE(pattern string, handler gin.HandlerFunc) {
	controller.r.DELETE(pattern, handler)
}
