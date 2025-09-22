package controller

import (
	"net/http"
	"orders/middleware"
	"orders/model"
	"orders/service"
	"strconv"

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
			return
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

	status, err := s.AddToCart(c.Request.Context(), input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	if !status {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "Failed to add to cart",
		})
		return
	}

	s.Commit()

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Product added to cart successfully",
	})
}

func UpdateCartItem(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "User is not logged in",
		})
		return
	}

	var input model.EditCartItem

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

	s.CartUpdateItem(c.Request.Context(), input)
	s.Commit()

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Cart updated successfully",
	})
}

type CheckoutInput struct {
	PaymentMethod string `json:"payment_method"`
	CartID        int    `json:"cart_id"`
	CartItemIDs   []int  `json:"cart_item_ids"`
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
	order, err := s.CreateOrder(c.Request.Context(), input.CartID, input.CartItemIDs, input.PaymentMethod)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	s.Commit()

	c.JSON(http.StatusOK, &model.OrderResponse{
		Success: true,
		Message: "Checked out successfully",
		Data:    []*model.Order{order},
	})
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

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "order ID is required",
		})
		return
	}

	orderID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "invalid order ID",
		})
		return
	}

	trackingInfo, err := s.OrderGetTrackingInfo(orderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.OrderTrackingResponse{
		Success: true,
		Message: "Order tracking info retrieved successfully",
		Data:    trackingInfo,
	})

}

type UpdateStatusInput struct {
	Status string `json:"status"`
}

func UpdateOrderStatus(c *gin.Context) {
	user := middleware.AuthContext(c.Request.Context())
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &model.GlobalResponse{
			Success: false,
			Message: "user not logged in",
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

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "order ID is required",
		})
		return
	}

	orderID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: "invalid order ID",
		})
		return
	}

	var input UpdateStatusInput

	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	status, err := s.OrderUpdateStatus(orderID, input.Status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if !status {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: "Failed to update order status",
		})
		return
	}

	s.Commit()

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Order status updated successfully",
	})
}
