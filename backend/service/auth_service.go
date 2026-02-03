package service

import (
	"exchangeapp/dao"
	"exchangeapp/models"
	"exchangeapp/utils"
)

type AuthService interface {
	Register(username, password string) (string, error)
	Login(username, password string) (string, error)
}

type authService struct {
	userDAO dao.UserDAO
}

func NewAuthService(userDAO dao.UserDAO) AuthService {
	return &authService{
		userDAO: userDAO,
	}
}

func (as *authService) Register(username, password string) (string, error) {
	// 检查用户是否已存在
	existingUser, _ := as.userDAO.GetUserByUsername(username)
	if existingUser != nil {
		return "", &AppError{Code: "USER_ALREADY_EXISTS", Message: "用户已存在"}
	}

	// 密码加密
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	// 生成Token
	token, err := utils.GenerateJWT(username)
	if err != nil {
		return "", err
	}

	// 创建用户
	user := &models.User{
		Username: username,
		Password: hashedPwd,
	}

	if err := as.userDAO.CreateUser(user); err != nil {
		return "", err
	}

	return token, nil
}

func (as *authService) Login(username, password string) (string, error) {
	// 查询用户
	user, err := as.userDAO.GetUserByUsername(username)
	if err != nil {
		return "", &AppError{Code: "USER_NOT_FOUND", Message: "用户不存在"}
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", &AppError{Code: "WRONG_PASSWORD", Message: "密码错误"}
	}

	// 生成Token
	token, err := utils.GenerateJWT(username)
	if err != nil {
		return "", err
	}

	return token, nil
}
