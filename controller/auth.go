package controller

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
	"github.com/rekib0023/event-horizon-gateway/utils"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	gRpc pb.AuthServiceClient
}

var authController *AuthController

func (controller *ControllerInterface) InitAuthController() {
	authController = &AuthController{
		gRpc: controller.gRpc,
	}

	POST("/auth/signup", authController.signup)
	POST("/auth/login", authController.login)
	GET("/auth/verify-token", authController.verifyToken)
	POST("/auth/refresh-token", authController.refreshToken)
}

func (o *AuthController) signup(c *gin.Context) {
	var reqData pb.SignupRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := o.gRpc.Signup(context.Background(), &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password})
	if err != nil {
		log.Printf("could not call Signup: %v", err)
		if s, ok := status.FromError(err); ok {
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.SetCookie("token", res.Token, 3600, "/", "", false, true)
	res.Token = ""
	c.JSON(http.StatusOK, res)
}

func (o *AuthController) login(c *gin.Context) {
	var reqData pb.LoginRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := o.gRpc.Login(context.Background(), &pb.LoginRequest{Email: reqData.Email, Password: reqData.Password})
	if err != nil {
		log.Printf("could not call Login: %v", err)
		if s, ok := status.FromError(err); ok {
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.SetCookie("token", res.Token, 3600, "/", "", false, true)
	res.Token = ""
	c.JSON(http.StatusOK, res)
}

func (o *AuthController) verifyToken(c *gin.Context) {
	var reqData pb.Token
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := o.gRpc.VerifyToken(context.Background(), &pb.Token{Token: reqData.Token})
	if err != nil {
		log.Printf("could not call VerifyToken: %v", err)
		if s, ok := status.FromError(err); ok {
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	c.JSON(http.StatusOK, res)
}

func (o *AuthController) refreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	token := parts[1]

	res, err := o.gRpc.RefreshToken(context.Background(), &pb.Token{Token: token})
	if err != nil {
		log.Printf("could not call RefreshToken: %v", err)
		if s, ok := status.FromError(err); ok {
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.SetCookie("token", res.Token, 3600, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{"message": "Token refreshed"})
}
