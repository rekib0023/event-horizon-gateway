package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
)

type AuthController struct {
}

var authController *AuthController

func (controller *Controller) InitAuthController() {
	authController = new(AuthController)

	POST("/auth/signup", authController.signup)
	POST("/auth/login", authController.login)
}

func (o *AuthController) signup(c *gin.Context) {
	var reqData pb.SignupRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := controller.gRpc.Signup(context.Background(), &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password})
	if err != nil {
		log.Printf("could not call Signup: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (o *AuthController) login(c *gin.Context) {
	var reqData pb.LoginRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := controller.gRpc.Login(context.Background(), &pb.LoginRequest{Email: reqData.Email, Password: reqData.Password})
	if err != nil {
		log.Printf("could not call Login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}
