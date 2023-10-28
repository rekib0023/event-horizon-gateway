package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
)

type ProfileController struct {
}

var profileController *ProfileController

func (controller *Controller) InitProfileController() {
	profileController = new(ProfileController)

	GET("/users", profileController.getUsers)
	GET("/users/:id", profileController.getUserById)
	PUT("/users/:id", profileController.updateUser)
	DELETE("/users/:id", profileController.deleteUser)
}

func (o *ProfileController) getUsers(c *gin.Context) {
	res, err := controller.gRpc.GetUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("could not call GetUsers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (o *ProfileController) getUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	res, err := controller.gRpc.GetUserById(context.Background(), &pb.UserId{Id: int32(id)})
	if err != nil {
		log.Printf("could not call GetUserById: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (o *ProfileController) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var reqData pb.SignupRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := controller.gRpc.UpdateUser(context.Background(), &pb.UpdateUserRequest{UserId: &pb.UserId{Id: int32(id)}, User: &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password}})
	if err != nil {
		log.Printf("could not call Update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (o *ProfileController) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	res, err := controller.gRpc.DeleteUser(context.Background(), &pb.UserId{Id: int32(id)})
	if err != nil {
		log.Printf("could not call Delete: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}