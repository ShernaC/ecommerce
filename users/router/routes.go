package router

import (
	"users/controller"
	"utils/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/approval", controller.Approval)

	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(), middleware.CORSMiddlewware(), middleware.IsLogin())
	{
		auth.POST("/refresh-token", controller.RefreshToken)
		auth.POST("/logout", controller.Logout)
		auth.GET("/profile/:id", controller.GetProfile)

		auth.POST("/seller/register", controller.RegisterSeller)
	}

	seller := r.Group("")
	seller.Use(middleware.AuthMiddleware(), middleware.CORSMiddlewware(), middleware.IsLogin(), middleware.IsSeller())
	{
		seller.GET("/seller/profile/:id", controller.SellerProfile)
	}

}
