package router

import (
	"users/controller"
	"users/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/approval", controller.Approval)

	auth := r.Group("")
	auth.Use(middleware.IsLogin())
	{
		auth.POST("/refresh-token", controller.RefreshToken)
		auth.POST("/logout", controller.Logout)
		auth.GET("/profile/:id", controller.GetProfile)

		auth.POST("/seller/register", controller.RegisterSeller)
		auth.GET("/seller/profile/:id", controller.SellerProfile)
	}
}
