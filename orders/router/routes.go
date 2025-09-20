package router

import (
	"orders/controller"

	"github.com/gin-gonic/gin"
)

// /cart → manage user’s shopping cart
// /checkout → protected endpoint to process a payment and create new order
// /orders → protected endpoint to retrieve a user’s order history
// /orders/{id}/track
// Interservice communication
// Product service → get product details and confirm inventory
// User service → confirm who place order

func ApiRouter(r *gin.Engine) {
	r.GET("/cart", controller.GetCart)
	r.POST("/cart", controller.AddToCart)
	r.POST("/cart/update", controller.UpdateCartItem)
	r.POST("/checkout", controller.Checkout)
	r.GET("/orders", controller.GetOrderHistory)
	r.GET("/orders/:id/track", controller.TrackOrder)
}
