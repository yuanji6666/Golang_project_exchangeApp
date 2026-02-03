package router

import ("github.com/gin-gonic/gin"
	"exchangeapp/controllers"
)
func registerArticlesRouter(r *gin.RouterGroup){
	r.POST("/articles", controllers.CreateArticle)
	r.GET("/articles", controllers.GetArticles )
	r.GET("/articles/:id", controllers.GetArticlesByID)
	r.POST("/articles/:id/like", controllers.LikeArticle)
	r.GET("/articles/:id/like", controllers.GetArticleLikes)
	
}