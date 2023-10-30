package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rekib0023/event-horizon-gateway/middlewares"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
)

type EventController struct {
	httpClient *http.Client
}

var eventController *EventController

func (controller *ControllerInterface) InitEventController() {
	eventController = &EventController{
		httpClient: &http.Client{},
	}

	controller.r.Use(middlewares.TokenAuthMiddleware(controller.gRpc))

	GET("/events", eventController.eventsPassThrough)
	POST("/events", eventController.eventsPassThrough)
	GET("/events/:eventId", eventController.eventsPassThrough)
	PUT("/events/:eventId", eventController.eventsPassThrough)
	DELETE("/events/:eventId", eventController.eventsPassThrough)
	GET("/events/:eventId/attendees", eventController.eventsPassThrough)
	POST("/events/:eventId/attendEvent", eventController.eventsPassThrough)
}

func (o *EventController) eventsPassThrough(c *gin.Context) {
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
			Errors string `json:"errors"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			o.jsonError(c, "Failed to decode error response", http.StatusInternalServerError)
			return
		}
		o.jsonError(c, errResp.Errors, resp.StatusCode)
		return
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		o.jsonError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (o *EventController) jsonError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}
