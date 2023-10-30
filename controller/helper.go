package controller

import (
	"github.com/gin-gonic/gin"
)

func POST(pattern string, handler gin.HandlerFunc) {
	controller.r.POST(pattern, handler)
}

func GET(pattern string, handler gin.HandlerFunc) {
	controller.r.GET(pattern, handler)
}

func PUT(pattern string, handler gin.HandlerFunc) {
	controller.r.PUT(pattern, handler)
}

func DELETE(pattern string, handler gin.HandlerFunc) {
	controller.r.DELETE(pattern, handler)
}

func ANY(pattern string, handler gin.HandlerFunc) {
	controller.r.Any(pattern, handler)
}

func USE(middlewares ...gin.HandlerFunc) {
	controller.r.Use(middlewares...)
}

func (o *ControllerInterface) jsonError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}
