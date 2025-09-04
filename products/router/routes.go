package router

import (
	"products/controller"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	r.POST("/product/create", controller.CreateProduct)
	r.GET("/product/:id", controller.ProductDetail)
	r.GET("/products/:seller_id}", controller.ProductList)
}
