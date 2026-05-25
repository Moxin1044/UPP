package service

import (
	"errors"

	"upp/internal/middleware"
	"upp/internal/model"
	"upp/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
	secret   string
	expire   int
}

func NewAuthService(userRepo *repository.UserRepository, secret string, expire int) *AuthService {
	return &AuthService{userRepo: userRepo, secret: secret, expire: expire}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo model.User  `json:"user_info"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role, s.secret, s.expire)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token, UserInfo: *user}, nil
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *AuthService) Register(req RegisterRequest) (*model.User, error) {
	_, err := s.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username: req.Username,
		Password: string(hashed),
		Role:     "user",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) InitAdmin() error {
	_, err := s.userRepo.FindByUsername("admin")
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := &model.User{
		Username: "admin",
		Password: string(hashed),
		Role:     "admin",
	}
	return s.userRepo.Create(admin)
}

func (s *AuthService) UserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}
