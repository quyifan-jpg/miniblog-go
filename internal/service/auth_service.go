package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"miniblog/internal/model"
	"miniblog/internal/pkg/jwt"
	"miniblog/internal/repository"
	"miniblog/internal/types"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register 用户注册
func (s *AuthService) Register(req *types.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	existUser, err := s.userRepo.GetByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existUser != nil {
		return nil, errors.New("username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
	}

	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("username or password incorrect")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username or password incorrect")
	}

	// 生成 Token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

