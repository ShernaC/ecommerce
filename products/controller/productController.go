package controller

import (
	"net/http"
	"products/model"
	"products/service"
	"utils/middleware"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not authenticated",
		})
		return
	}

	isSeller, err := service.CheckSellerExists(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "failed to check user permissions",
		})
		return
	}

	if !isSeller {
		c.AbortWithStatusJSON(http.StatusForbidden, &model.GlobalResponse{
			Success: false,
			Message: "user does not have seller permissions",
		})
		return
	}

	var input model.NewProduct

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	input.SellerID = user.ID

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

	s.ProductCreate(c.Request.Context(), input)
	s.Commit()

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Product successfully created",
	})
}

func ProductDetail(c *gin.Context) {
	var input int

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	s := service.GetService()
	defer func() {
		if r := recover(); r != nil {
			err := s.ErrorCheck(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}()

	s.ProductGetByID(c.Request.Context(), input)

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Product detail retrieved successfully",
	})
}

func ProductList(c *gin.Context) {
	var input int

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	s := service.GetService()
	defer func() {
		if r := recover(); r != nil {
			err := s.ErrorCheck(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}()

	s.ProductGetProductsBySellerID(c.Request.Context(), input)

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Seller's product details retrieved successfully",
	})
}
