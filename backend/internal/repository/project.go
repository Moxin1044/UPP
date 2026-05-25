package repository

import (
	"upp/internal/model"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) FindByUserID(userID uint) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) FindByID(id uint) (*model.Project, error) {
	var project model.Project
	if err := r.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&model.Project{}, id).Error
}
