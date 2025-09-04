package middleware

import (
	"context"
	"net/http"
	"products/model"
	"products/tools"
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
	// return a function with gin context
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			c.Next()
			return
		}

		authTokens := strings.Split(authToken, " ")
		if len(authTokens) != 2 || authTokens[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		// 3. Parse and validate token
		token, err := tools.ValidateToken(authTokens[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
				Success: false,
				Message: "invalid token",
			})
			return
		}

		// 4. Extract and validate claims
		claims, ok := token.Claims.(*tools.Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
				Success: false,
				Message: "invalid token claims",
			})
			return
		}

		// 5. Check token expiration
		if exp := claims.ExpiresAt; exp != 0 {
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
