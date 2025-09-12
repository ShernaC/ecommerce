package controller

import (
	"net/http"
	"orders/model"
	"orders/service"
	"users/middleware"

	"github.com/gin-gonic/gin"
)

func GetCart(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not logged in",
		})
		return
	}

	s := service.GetService()
	defer func() {
		r := recover()
		if r != nil {
			err := s.ErrorCheck(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
		}
	}()

	cart, err := s.CartGetDetails(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, &model.CartResponse{
		Success: true,
		Message: "Cart details retrieved successfully",
		Data:    *cart,
	})
}

func AddToCart(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not logged in",
		})
		return
	}

	var input model.CartItemInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	s := service.GetTransaction()
	defer func() {
		r := recover()
		if r != nil {
			err := s.Rollback(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}()

	s.AddToCart(c.Request.Context(), input)

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Product added to cart successfully",
	})
}

func Checkout(c *gin.Context) {

}

func GetOrderHistory(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not logged in",
		})
		return
	}

	s := service.GetService()
	defer func() {
		r := recover()
		if r != nil {
			err := s.ErrorCheck(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
		}
	}()

	orders, err := s.OrderGetHistoryByUserID(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, &model.OrderResponse{
		Success: true,
		Message: "Order history retrieved successfully",
		Data:    orders,
	})
}

func TrackOrder(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not logged in",
		})
		return
	}

	s := service.GetService()
	defer func() {
		r := recover()
		if r != nil {
			err := s.ErrorCheck(r)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
				Success: false,
				Message: err.Error(),
			})
		}
	}()

}
