package controllers

import (
	"exchangeapp/dao"
	"exchangeapp/models"
	"exchangeapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	articleService service.ArticleService
	likeService    service.LikeService
}

func NewArticleController() *ArticleController {
	articleDAO := dao.NewArticleDAO()
	likeDAO := dao.NewLikeDAO()
	articleService := service.NewArticleService(articleDAO, likeDAO)
	likeService := service.NewLikeService(likeDAO)

	return &ArticleController{
		articleService: articleService,
		likeService:    likeService,
	}
}

func (ac *ArticleController) CreateArticle(ctx *gin.Context) {
	var article models.Article

	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := ac.articleService.CreateArticle(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func (ac *ArticleController) GetArticles(ctx *gin.Context) {
	articles, err := ac.articleService.GetArticles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

func (ac *ArticleController) GetArticleByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	article, err := ac.articleService.GetArticleByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (ac *ArticleController) LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")

	if err := ac.likeService.LikeArticle(articleID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "成功点赞"})
}

func (ac *ArticleController) GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likes, err := ac.likeService.GetArticleLikes(articleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}
