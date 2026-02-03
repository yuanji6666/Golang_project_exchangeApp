package service

import (
	"exchangeapp/dao"
	"strconv"
)

type LikeService interface {
	LikeArticle(articleID string) error
	GetArticleLikes(articleID string) (int64, error)
	UnlikeArticle(articleID string) error
}

type likeService struct {
	likeDAO dao.LikeDAO
}

func NewLikeService(likeDAO dao.LikeDAO) LikeService {
	return &likeService{
		likeDAO: likeDAO,
	}
}

func (ls *likeService) LikeArticle(articleID string) error {
	// 验证articleID
	if _, err := strconv.ParseUint(articleID, 10, 64); err != nil {
		return &AppError{Code: "INVALID_ARTICLE_ID", Message: "无效的文章ID"}
	}

	return ls.likeDAO.LikeArticle(articleID)
}

func (ls *likeService) GetArticleLikes(articleID string) (int64, error) {
	// 验证articleID
	if _, err := strconv.ParseUint(articleID, 10, 64); err != nil {
		return 0, &AppError{Code: "INVALID_ARTICLE_ID", Message: "无效的文章ID"}
	}

	return ls.likeDAO.GetArticleLikes(articleID)
}

func (ls *likeService) UnlikeArticle(articleID string) error {
	// 验证articleID
	if _, err := strconv.ParseUint(articleID, 10, 64); err != nil {
		return &AppError{Code: "INVALID_ARTICLE_ID", Message: "无效的文章ID"}
	}

	return ls.likeDAO.UnlikeArticle(articleID)
}
