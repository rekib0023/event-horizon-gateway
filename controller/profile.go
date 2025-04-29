package controller

import (
	"context"
	"fmt" // Added for formatted logging
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rekib0023/event-horizon-gateway/middlewares"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
	"github.com/rekib0023/event-horizon-gateway/utils"
	"google.golang.org/grpc/status"
)

// ProfileController handles user profile related API requests.
type ProfileController struct {
	gRpc pb.AuthServiceClient
}

// profileController is a singleton instance of ProfileController.
var profileController *ProfileController

// InitProfileController initializes the ProfileController and its routes.
// It also applies the token authentication middleware to these routes.
func (c *ControllerInterface) InitProfileController() {
	profileController = &ProfileController{
		gRpc: c.gRpc,
	}

	// Apply the token authentication middleware to all routes defined in this controller.
	c.r.Use(middlewares.TokenAuthMiddleware(c.gRpc))

	// Define the API endpoints for user profiles.
	GET("/users", profileController.getUsers)
	GET("/users/:userId", profileController.getUserById)
	PUT("/users/:userId", profileController.updateUser)
	DELETE("/users/:userId", profileController.deleteUser)
	GET("/users/:userId/events", c.eventsPassThrough) // This route seems to pass through to another service.
}

// Helper function to parse and validate userId
func parseAndValidateUserId(c *gin.Context) (int32, error) {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("could not parse userId: %v", err)
		return 0, fmt.Errorf("invalid user ID format")
	}
	if id <= 0 {
		log.Printf("invalid userId: %d", id)
		return 0, fmt.Errorf("user ID must be a positive integer")
	}
	return int32(id), nil
}

// getUsers handles the request to retrieve all users.
func (o *ProfileController) getUsers(c *gin.Context) {
	// Check if the user information is available in the context after authentication.
	_, exists := c.Get("user")
	if !exists {
		log.Println("user information not found in context") // More specific log message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: User context missing"}) // More informative error
		return
	}

	// Call the GetUsers gRPC method of the Auth service.
	res, err := o.gRpc.GetUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("could not call GetUsers: %v", err)
		// Check if the error is a gRPC status error.
		if s, ok := status.FromError(err); ok {
			// Convert the gRPC status code to an HTTP status code.
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			// Handle non-gRPC errors.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: Failed to retrieve users"}) // More specific error
			return
		}
	}
	// Respond with the list of users in JSON format.
	c.JSON(http.StatusOK, res)
}

// getUserById handles the request to retrieve a specific user by their ID.
func (o *ProfileController) getUserById(c *gin.Context) {
	// Use the helper function to parse and validate the userId.
	userId, err := parseAndValidateUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the GetUserById gRPC method with the extracted user ID.
	res, err := o.gRpc.GetUserById(context.Background(), &pb.UserId{Id: userId})
	if err != nil {
		log.Printf("could not call GetUserById for ID %d: %v", userId, err) // Include the ID in the log
		// Check if the error is a gRPC status error.
		if s, ok := status.FromError(err); ok {
			// Convert the gRPC status code to an HTTP status code.
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			// Handle non-gRPC errors.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: Failed to retrieve user"}) // More specific error
			return
		}
	}
	// Respond with the user information in JSON format.
	c.JSON(http.StatusOK, res)
}

// updateUser handles the request to update an existing user's information.
func (o *ProfileController) updateUser(c *gin.Context) {
	// Use the helper function to parse and validate the userId.
	userId, err := parseAndValidateUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define a struct to hold the request data for updating a user.
	var reqData pb.SignupRequest
	// Bind the JSON request body to the reqData struct.
	if err := c.ShouldBindJSON(&reqData); err != nil {
		log.Printf("invalid request body for update user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"}) // More specific error
		return
	}

	// Basic validation for required fields (as a fresher might implement)
	if reqData.FirstName == "" || reqData.LastName == "" || reqData.UserName == "" || reqData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields in request body"})
		return
	}

	// Call the UpdateUser gRPC method with the user ID and updated information.
	res, err := o.gRpc.UpdateUser(context.Background(), &pb.UpdateUserRequest{UserId: &pb.UserId{Id: userId}, User: &pb.SignupRequest{FirstName: reqData.FirstName, LastName: reqData.LastName, UserName: reqData.UserName, Email: reqData.Email, Password: reqData.Password}})
	if err != nil {
		log.Printf("could not call UpdateUser for ID %d: %v", userId, err) // Include the ID in the log
		// Check if the error is a gRPC status error.
		if s, ok := status.FromError(err); ok {
			// Convert the gRPC status code to an HTTP status code.
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			// Handle non-gRPC errors.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: Failed to update user"}) // More specific error
			return
		}
	}
	// Respond with the updated user information in JSON format.
	c.JSON(http.StatusOK, res)
}

// deleteUser handles the request to delete a user by their ID.
func (o *ProfileController) deleteUser(c *gin.Context) {
	// Use the helper function to parse and validate the userId.
	userId, err := parseAndValidateUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the DeleteUser gRPC method with the user ID to be deleted.
	res, err := o.gRpc.DeleteUser(context.Background(), &pb.UserId{Id: userId})
	if err != nil {
		log.Printf("could not call DeleteUser for ID %d: %v", userId, err) // Include the ID in the log
		// Check if the error is a gRPC status error.
		if s, ok := status.FromError(err); ok {
			// Convert the gRPC status code to an HTTP status code.
			httpStatusCode := utils.GetHttpStatusCode(s.Code())
			c.JSON(httpStatusCode, gin.H{"error": s.Message()})
			return
		} else {
			// Handle non-gRPC errors.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: Failed to delete user"}) // More specific error
			return
		}
	}
	// Respond with a 204 No Content status to indicate successful deletion.
	c.JSON(http.StatusNoContent, res)
}
