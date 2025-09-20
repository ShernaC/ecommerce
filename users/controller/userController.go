package controller

import (
	"net/http"
	"strconv"
	"time"
	grpcclient "users/grpc_client"
	"users/model"
	"users/service"
	"utils/orders"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input model.NewUser

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	s := service.GetTransaction()
	defer func() {
		if r := recover(); r != nil {
			err := s.Rollback(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}()

	user, err := s.UserRegister(c.Request.Context(), input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// create cart here
	success, err := grpcclient.CreateCart(c.Request.Context(), &orders.CreateCartRequest{UserId: int64(user.ID)})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	if !success.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "failed to create cart for user",
		})
		return
	}

	s.Commit()

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "User successfully created",
	})
}

func Login(c *gin.Context) {
	var input model.UserLogin

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	s := service.GetTransaction()
	defer func() {
		if r := recover(); r != nil {
			err := s.Rollback(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}()

	resp, _ := s.UserLogin(c.Request.Context(), input)
	s.Commit()

	c.SetCookie("refresh_token", resp.Data[1].Token, int(24*time.Hour.Seconds()), "/", "localhost", true, true)

	c.JSON(http.StatusOK, &model.UserLoginResponse{
		Success: true,
		Message: "Login successful",
		Data: []*model.UserLoginResponseNode{
			{
				TokenType: resp.Data[0].TokenType,
				Token:     resp.Data[0].Token,
				UserData:  resp.Data[0].UserData,
			},
		},
	})
}

func Logout(c *gin.Context) {
	_, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "You have logged out. Please log in again",
		})
	}

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

	err = s.UserRevokeToken(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Logout successful",
	})
}

func GetProfile(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "Invalid user id",
		})
		return
	}

	s := service.GetService()
	defer func() {
		if r := recover(); r != nil {
			err := s.ErrorCheck()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}
		}
	}()

	user, err := s.UserGetByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "Error receiving user info",
		})
	}

	c.JSON(http.StatusOK, &model.UserResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data: model.UserData{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
	})
}
