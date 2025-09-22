package router

import (
	"orders/controller"
	"utils/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(), middleware.IsLogin())
	{
		auth.GET("/cart", controller.GetCart)
		auth.POST("/cart", controller.AddToCart)
		auth.POST("/cart/update", controller.UpdateCartItem)
		auth.POST("/checkout", controller.Checkout)
		auth.GET("/orders", controller.GetOrderHistory)
		auth.GET("/orders/:id/track", controller.TrackOrder)
	}

	seller := r.Group("")
	seller.Use(middleware.AuthMiddleware(), middleware.CORSMiddlewware(), middleware.IsLogin(), middleware.IsSeller())
	{
		seller.POST("/orders/:id/updateStatus", controller.UpdateOrderStatus)
	}
}
