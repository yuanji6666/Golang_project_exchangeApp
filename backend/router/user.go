package router

import (
	"exchangeapp/controllers"

	"github.com/gin-gonic/gin"
)

func registerUserRouter(r *gin.RouterGroup) {
	ac := controllers.NewAuthController()

	r.POST("/login", ac.Login)
	r.POST("/register", ac.Register)
}
