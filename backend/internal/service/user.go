package service

import (
	"upp/internal/model"
	"upp/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListAll() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

func (s *UserService) Create(req CreateUserRequest) (*model.User, error) {
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, ErrUsernameExists
	}
	hashed, err := bcryptHash(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username: req.Username,
		Password: hashed,
		Role:     req.Role,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

func (s *UserService) UpdatePassword(id uint, req UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	hashed, err := bcryptHash(req.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
