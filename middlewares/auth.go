package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/rekib0023/event-horizon-gateway/proto"
)

type AuthMiddleware struct {
}

var authMiddleware *AuthMiddleware

func TokenAuthMiddleware(gRpc pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		res, err := gRpc.VerifyToken(context.Background(), &pb.Token{Token: token})
		if err != nil {
			log.Printf("could not call VerifyToken: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}
		c.Set("user", res)
		c.Next()
	}
}
