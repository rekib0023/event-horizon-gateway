package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
)

func RegisterAuthEndpoints(r *gin.RouterGroup, gRpc pb.AuthServiceClient) {
	r.POST("/signup", func(c *gin.Context) {
		var reqData pb.SignupRequest
		if err := c.ShouldBindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		res, err := gRpc.Signup(context.Background(), &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password})
		if err != nil {
			log.Printf("could not call Signup: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.POST("/login", func(c *gin.Context) {
		var reqData pb.LoginRequest
		if err := c.ShouldBindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		res, err := gRpc.Login(context.Background(), &pb.LoginRequest{Email: reqData.Email, Password: reqData.Password})
		if err != nil {
			log.Printf("could not call Login: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
