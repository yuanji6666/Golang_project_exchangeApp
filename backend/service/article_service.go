package service

import (
	"encoding/json"
	"exchangeapp/dao"
	"exchangeapp/global"
	"exchangeapp/models"
	"time"

	"github.com/go-redis/redis"
)

type ArticleService interface {
	CreateArticle(article *models.Article) error
	GetArticles() ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	UpdateArticle(article *models.Article) error
	DeleteArticle(id uint) error
}

type articleService struct {
	articleDAO dao.ArticleDAO
	likeDAO    dao.LikeDAO
}

const articleCacheKey = "articles"
const articleCacheTTL = 10 * time.Minute

func NewArticleService(articleDAO dao.ArticleDAO, likeDAO dao.LikeDAO) ArticleService {
	return &articleService{
		articleDAO: articleDAO,
		likeDAO:    likeDAO,
	}
}

func (as *articleService) CreateArticle(article *models.Article) error {
	if err := as.articleDAO.CreateArticle(article); err != nil {
		return err
	}

	// 删除缓存
	_ = global.RedisDb.Del(articleCacheKey).Err()
	return nil
}

func (as *articleService) GetArticles() ([]models.Article, error) {
	// 尝试从缓存获取
	cachedData, err := global.RedisDb.Get(articleCacheKey).Result()

	if err == nil {
		// 缓存命中
		var articles []models.Article
		if err := json.Unmarshal([]byte(cachedData), &articles); err == nil {
			return articles, nil
		}
	} else if err != redis.Nil {
		// Redis错误，继续执行数据库查询
		return nil, err
	}

	// 缓存未命中或过期，查询数据库
	articles, err := as.articleDAO.GetArticles()
	if err != nil {
		return nil, err
	}

	// 将数据写入缓存
	articlesJSON, _ := json.Marshal(articles)
	_ = global.RedisDb.Set(articleCacheKey, articlesJSON, articleCacheTTL).Err()

	return articles, nil
}

func (as *articleService) GetArticleByID(id uint) (*models.Article, error) {
	article, err := as.articleDAO.GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	// 获取点赞数
	likes, _ := as.likeDAO.GetArticleLikes("")
	if article != nil {
		article.Likes = int(likes)
	}

	return article, nil
}

func (as *articleService) UpdateArticle(article *models.Article) error {
	if err := as.articleDAO.UpdateArticle(article); err != nil {
		return err
	}

	// 删除缓存
	_ = global.RedisDb.Del(articleCacheKey).Err()
	return nil
}

func (as *articleService) DeleteArticle(id uint) error {
	if err := as.articleDAO.DeleteArticle(id); err != nil {
		return err
	}

	// 删除缓存
	_ = global.RedisDb.Del(articleCacheKey).Err()
	return nil
}
