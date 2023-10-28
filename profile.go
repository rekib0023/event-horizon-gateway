// profile.go
package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
)

func RegisterProfileEndpoints(r *gin.Engine, gRpc pb.AuthServiceClient) {
	r.GET("/api/users", func(c *gin.Context) {
		res, err := gRpc.GetUsers(context.Background(), &pb.Empty{})
		if err != nil {
			log.Printf("could not call GetUsers: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("could not parse id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		res, err := gRpc.GetUserById(context.Background(), &pb.UserId{Id: int32(id)})
		if err != nil {
			log.Printf("could not call GetUserById: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.PUT("/api/users/:id", func(c *gin.Context) {
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

		res, err := gRpc.UpdateUser(context.Background(), &pb.UpdateUserRequest{UserId: &pb.UserId{Id: int32(id)}, User: &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password}})
		if err != nil {
			log.Printf("could not call Update: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.DELETE("/api/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("could not parse id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		res, err := gRpc.DeleteUser(context.Background(), &pb.UserId{Id: int32(id)})
		if err != nil {
			log.Printf("could not call Delete: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
