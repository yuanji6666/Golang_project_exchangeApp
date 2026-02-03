package dao

import (
	"exchangeapp/global"
	"exchangeapp/models"

	"gorm.io/gorm"
)

type ArticleDAO interface {
	CreateArticle(article *models.Article) error
	GetArticles() ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	UpdateArticle(article *models.Article) error
	DeleteArticle(id uint) error
	GetAllArticlesCount() (int64, error)
}

type articleDAO struct {
	db *gorm.DB
}

func NewArticleDAO() ArticleDAO {
	return &articleDAO{
		db: global.Db,
	}
}

func (ad *articleDAO) CreateArticle(article *models.Article) error {
	return ad.db.Create(article).Error
}

func (ad *articleDAO) GetArticles() ([]models.Article, error) {
	var articles []models.Article
	if err := ad.db.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (ad *articleDAO) GetArticleByID(id uint) (*models.Article, error) {
	var article models.Article
	if err := ad.db.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (ad *articleDAO) UpdateArticle(article *models.Article) error {
	return ad.db.Save(article).Error
}

func (ad *articleDAO) DeleteArticle(id uint) error {
	return ad.db.Delete(&models.Article{}, id).Error
}

func (ad *articleDAO) GetAllArticlesCount() (int64, error) {
	var count int64
	if err := ad.db.Model(&models.Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
