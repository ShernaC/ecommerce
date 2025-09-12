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

	s.CartGetDetails(c.Request.Context())

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Cart details retrieved successfully",
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

type CheckoutInput struct {
	PaymentMethod string     `json:"payment_method"`
	Cart          model.Cart `json:"cart"`
}

func Checkout(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not logged in",
		})
		return
	}

	var input CheckoutInput

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

	// Only create order once payment is successful

	s.CreateOrder(c.Request.Context(), input.Cart, input.PaymentMethod)
}

func GetOrderHistory(c *gin.Context) {

}

func TrackOrder(c *gin.Context) {

}
