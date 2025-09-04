package controller

import (
	"fmt"
	"net/http"
	"time"
	"users/model"
	"users/service"
	"users/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RefreshToken(c *gin.Context) {
	// Get refresh token from HttpOnly cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "Unauthorized - refresh token not found",
		})
		return
	}

	// validate refresh token's signature and expiration
	jwtToken, err := tools.ValidateToken(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "Unauthorized - invalid refresh token",
		})
		return
	}

	// extract claims and jti from token
	claims, ok := jwtToken.Claims.(*tools.Claims)
	if !ok || claims.Id == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "Unauthorized - invalid token claims or missing session ID",
		})
		return
	}

	// check jti
	s := service.GetService()
	defer func() {
		if r := recover(); r != nil {
			err := s.ErrorCheck(r)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}
		}
	}()

	isValid, err := s.UserCheckTokenValid(c.Request.Context(), claims.ID, claims.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	if !isValid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "invalid refresh token",
		})
		return
	}

	// revoke old jti -> prevent replay attacks
	err = s.UserRevokeToken(c.Request.Context())
	if err != nil {
		fmt.Printf("Warning: failed to revoke old refresh token JTI: %v", err)
	}

	// new jti and refresh token
	jti := uuid.New().String()

	newRefreshToken, err := tools.CreateToken(claims.ID, claims.Email, claims.Role, 24*time.Hour, jti)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// store jti in db
	user, err := s.UserUpdateRememberToken(c.Request.Context(), claims.ID, jti)
	if err != nil {
		panic(fmt.Errorf("failed to store refresh token jti: %v", err))
	}

	// set new refresh token as cookie
	c.SetCookie("refresh_token", newRefreshToken, 3600*24, "/", "", true, true)

	// generate new access token
	newAccessToken, err := tools.CreateToken(claims.ID, claims.Email, claims.Role, 30*time.Minute, jti)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.UserLoginResponse{
		Success: true,
		Message: "access token refreshed successfully",
		Data: []*model.UserLoginResponseNode{
			{
				TokenType: "access",
				Token:     newAccessToken,
				UserData: model.UserData{
					ID:        user.ID,
					Email:     user.Email,
					Name:      user.Name,
					Phone:     user.Phone,
					CreatedAt: user.CreatedAt,
				},
			},
		},
	})

}
