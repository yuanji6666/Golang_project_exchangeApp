package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET","POST"},
		AllowHeaders:     []string{"Origin","Content-Type","Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
  }))
	auth := r.Group("/api/auth") 
	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register)

	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRates)
	api.Use(middlewares.AuthMiddleWare())
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles )
		api.GET("/articles/:id", controllers.GetArticlesByID)

		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
	}
	return r
	
}