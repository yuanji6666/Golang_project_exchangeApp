package dao

import (
	"exchangeapp/global"
	"exchangeapp/models"

	"gorm.io/gorm"
)

type UserDAO interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userDAO struct {
	db *gorm.DB
}

func NewUserDAO() UserDAO {
	return &userDAO{
		db: global.Db,
	}
}

func (ud *userDAO) CreateUser(user *models.User) error {
	return ud.db.Create(user).Error
}

func (ud *userDAO) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := ud.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ud *userDAO) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := ud.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ud *userDAO) UpdateUser(user *models.User) error {
	return ud.db.Save(user).Error
}

func (ud *userDAO) DeleteUser(id uint) error {
	return ud.db.Delete(&models.User{}, id).Error
}
