package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
	"google.golang.org/grpc"
)

type ControllerInterface struct {
	r          *gin.RouterGroup
	gRpc       pb.AuthServiceClient
	httpClient *http.Client
}

var controller *ControllerInterface

func Init() {
	controller.InitAuthController()
	controller.InitProfileController()
	controller.InitEventController()
}

var e *gin.Engine

func Start() {
	e = gin.Default()

	conn, err := grpc.Dial(os.Getenv("AUTH_SVC"), grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
	} else {
		defer conn.Close()
		gRpc := pb.NewAuthServiceClient(conn)
		apiGroup := e.Group("/api")
		controller = &ControllerInterface{
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

func (o *ControllerInterface) eventsPassThrough(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		o.jsonError(c, "Internal server error", http.StatusInternalServerError)
		return
	}

	currentUser, ok := userValue.(*pb.TokenVerification)
	if !ok {
		o.jsonError(c, "Internal server error", http.StatusInternalServerError)
		return
	}

	endpoint := strings.TrimPrefix(c.Request.URL.Path, "/api")
	url := os.Getenv("EVENT_MGT_SVC") + endpoint
	req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)

	if err != nil {
		o.jsonError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-User-ID", currentUser.Id)
	req.Header.Set("X-User-Email", currentUser.Email)
	req.Header.Add("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		o.jsonError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			o.jsonError(c, "Failed to decode error response", http.StatusInternalServerError)
			return
		}
		var errMsgs []string
		for _, e := range errResp.Errors {
			errMsgs = append(errMsgs, e.Message)
		}
		errMsg := strings.Join(errMsgs, ", ")
		o.jsonError(c, errMsg, resp.StatusCode)
		return
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		o.jsonError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}
