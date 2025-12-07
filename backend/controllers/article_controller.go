package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateArticle(ctx *gin.Context){
	var article models.Article

	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error" : err ,
		})
		return
	}

	if err := global.Db.AutoMigrate(&article);err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error" : err ,
		})
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error" : err ,
		})
		return
	}

	ctx.JSON(http.StatusCreated, article )
}

func GetArticles(ctx *gin.Context){
	var articles []models.Article

	if err := global.Db.Find(&articles).Error; err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error" : err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK,articles)
}

func GetArticlesByID(ctx *gin.Context){
	id :=ctx.Param("id")
	var article models.Article

	if err := global.Db.Where("id=?",id).First(&article).Error; err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}else{
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"error" : err.Error(),
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, article)
}
