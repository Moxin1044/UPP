package service

import (
	"upp/internal/model"
	"upp/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
}

func (s *ProjectService) Create(userID uint, req CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) ListByUser(userID uint) ([]model.Project, error) {
	return s.repo.FindByUserID(userID)
}

func (s *ProjectService) GetByID(id uint) (*model.Project, error) {
	return s.repo.FindByID(id)
}

func (s *ProjectService) Update(id uint, req UpdateProjectRequest) (*model.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	project.Name = req.Name
	project.Description = req.Description
	if err := s.repo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}
