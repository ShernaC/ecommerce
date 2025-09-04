package controller

import (
	"net/http"
	"strconv"
	"users/model"
	"users/service"

	"github.com/gin-gonic/gin"
)

func RegisterSeller(c *gin.Context) {
	var input model.NewSeller

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

	s.SellerRegister(c.Request.Context(), input)
	s.Commit()

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Seller successfully created",
	})
}

func SellerProfile(c *gin.Context) {
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

	seller, err := s.SellerGetByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "Error receiving seller info",
		})
	}

	if seller.IsApproved == string(service.APPROVAL_TYPE_REJECTED) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "Seller account is not approved",
		})
		return
	} else if seller.IsApproved == string(service.APPROVAL_TYPE_PENDING) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "Seller account is pending approval",
		})
		return
	}

	userData, err := s.UserGetByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "Error receiving seller's general user info",
		})
	}

	c.JSON(http.StatusOK, &model.SellerResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data: model.SellerData{
			UserData: model.UserData{
				ID:    id,
				Name:  userData.Name,
				Email: userData.Email,
				Phone: userData.Phone,
			},
			BusinessName: seller.BusinessName,
			Address:      seller.Address,
			IsApproved:   seller.IsApproved,
			CreatedAt:    seller.CreatedAt,
		},
	})
}
