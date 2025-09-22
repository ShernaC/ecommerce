package router

import (
	"products/controller"
	"utils/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	r.POST("/product/create", controller.CreateProduct)
	r.GET("/product/:id", controller.ProductDetail)

	seller := r.Group("")
	seller.Use(middleware.AuthMiddleware(), middleware.CORSMiddlewware(), middleware.IsLogin(), middleware.IsSeller())
	{
		seller.GET("/products/:seller_id}", controller.ProductList)
	}
}
