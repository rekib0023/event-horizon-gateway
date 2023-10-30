package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekib0023/event-horizon-gateway/middlewares"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
	"github.com/rekib0023/event-horizon-gateway/utils"
	"google.golang.org/grpc/status"
)

type ProfileController struct {
	gRpc pb.AuthServiceClient
}

var profileController *ProfileController

func (c *ControllerInterface) InitProfileController() {
	profileController = &ProfileController{
		gRpc: c.gRpc,
	}

	c.r.Use(middlewares.TokenAuthMiddleware(c.gRpc))

	GET("/users", profileController.getUsers)
	GET("/users/:userId", profileController.getUserById)
	PUT("/users/:userId", profileController.updateUser)
	DELETE("/users/:userId", profileController.deleteUser)
}

func (o *ProfileController) getUsers(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	res, err := o.gRpc.GetUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("could not call GetUsers: %v", err)
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

func (o *ProfileController) getUserById(c *gin.Context) {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse userId: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	res, err := o.gRpc.GetUserById(context.Background(), &pb.UserId{Id: int32(id)})
	if err != nil {
		log.Printf("could not call GetUserById: %v", err)
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

func (o *ProfileController) updateUser(c *gin.Context) {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse userId: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var reqData pb.SignupRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := o.gRpc.UpdateUser(context.Background(), &pb.UpdateUserRequest{UserId: &pb.UserId{Id: int32(id)}, User: &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password}})
	if err != nil {
		log.Printf("could not call Update: %v", err)
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

func (o *ProfileController) deleteUser(c *gin.Context) {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse userId: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	res, err := o.gRpc.DeleteUser(context.Background(), &pb.UserId{Id: int32(id)})
	if err != nil {
		log.Printf("could not call Delete: %v", err)
		if s, ok := status.FromError(err); ok {
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	c.JSON(http.StatusNoContent, res)
}
