package router

import (
	"exchangeapp/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://www.yuanji666.fun"},
		AllowMethods:     []string{"GET","POST"},
		AllowHeaders:     []string{"Origin","Content-Type","Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
  	}))
	auth := r.Group("/api/auth") 
	registerUserRouter(auth)

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleWare())
	{
		registerExchangeRatesRouter(api)
		registerArticlesRouter(api)
	}
	return r
	
}