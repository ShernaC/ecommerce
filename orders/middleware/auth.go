package middleware

import (
	"context"
	"log"
	"net/http"
	"orders/model"
	"orders/tools"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var CtxKey = &contextKey{}

type contextKey struct {
	name string
}

type User struct {
	ID int `json:"id"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			c.Next()
			return
		}

		authTokens := strings.Split(authToken, " ")
		if authTokens == nil || authTokens[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
				Success: false,
				Message: "invalid authorisation token",
			})
			return
		}

		token, err := tools.ValidateToken(authTokens[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: "error validating token",
			})
			return
		}

		claims, ok := token.Claims.(*tools.Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: "error extracting token claims",
			})
			return
		}

		exp := claims.ExpiresAt
		if exp != 0 {
			if time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
					Success: false,
					Message: "token expired",
				})
				return
			}
		}

		ctx := context.WithValue(c.Request.Context(), CtxKey, &User{
			ID: claims.ID,
		})

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AuthContext(ctx context.Context) *User {
	raw, _ := ctx.Value(CtxKey).(*User)
	return raw
}

func IsLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if AuthContext(c.Request.Context()) == nil {
			log.Println("No context found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
				Success: false,
				Message: "Invalid token",
			})
			return
		}
		c.Next()
	}
}
