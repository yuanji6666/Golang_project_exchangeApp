package router

import (
	"exchangeapp/controllers"

	"github.com/gin-gonic/gin"
)

func registerArticlesRouter(r *gin.RouterGroup) {
	ac := controllers.NewArticleController()

	r.POST("/articles", ac.CreateArticle)
	r.GET("/articles", ac.GetArticles)
	r.GET("/articles/:id", ac.GetArticleByID)
	r.POST("/articles/:id/like", ac.LikeArticle)
	r.GET("/articles/:id/like", ac.GetArticleLikes)
}
