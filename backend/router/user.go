package router

import ("github.com/gin-gonic/gin"
	"exchangeapp/controllers"
)

func registerUserRouter(r *gin.RouterGroup){
	
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
}