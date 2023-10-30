package controller

import (
	"net/http"

	"github.com/rekib0023/event-horizon-gateway/middlewares"
)

func (o *ControllerInterface) InitEventController() {
	o.httpClient = &http.Client{}

	USE(middlewares.TokenAuthMiddleware(o.gRpc))

	GET("/events", o.eventsPassThrough)
	POST("/events", o.eventsPassThrough)
	GET("/events/:eventId", o.eventsPassThrough)
	PUT("/events/:eventId", o.eventsPassThrough)
	DELETE("/events/:eventId", o.eventsPassThrough)
	GET("/events/:eventId/attendees", o.eventsPassThrough)
	POST("/events/:eventId/attendEvent", o.eventsPassThrough)
	POST("/events/:eventId/register", o.eventsPassThrough)
	GET("/users/:userId/events", o.eventsPassThrough)
}
