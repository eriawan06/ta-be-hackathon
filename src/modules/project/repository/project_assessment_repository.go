package repository

import (
	"be-sagara-hackathon/src/modules/project/model"
	"be-sagara-hackathon/src/utils/constants"
	"gorm.io/gorm"
)

type ProjectAssessmentRepository interface {
	Create(assessment model.ProjectAssessment) error
	CreateBatch(assessments []model.ProjectAssessment) error
	FindByProjectID(projectID uint) (assessments []model.ProjectAssessment, err error)
	FindByProjectIDAndJudgeID(projectID, judgeID uint) (assessment []model.ProjectAssessment, err error)
}

type ProjectAssessmentRepositoryImpl struct {
	DB *gorm.DB
}

func NewProjectAssessmentRepository(db *gorm.DB) ProjectAssessmentRepository {
	return &ProjectAssessmentRepositoryImpl{DB: db}
}

func (repository *ProjectAssessmentRepositoryImpl) Create(assessment model.ProjectAssessment) error {
	if err := repository.DB.Create(&assessment).Error; err != nil {
		return err
	}
	return nil
}

func (repository *ProjectAssessmentRepositoryImpl) CreateBatch(assessments []model.ProjectAssessment) error {
	tx := repository.DB.Begin()
	if err := tx.Create(&assessments).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Project{}).
		Where("id=?", assessments[0].ProjectID).
		Update("status", constants.ProjectStatusAssessed).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repository *ProjectAssessmentRepositoryImpl) FindByProjectID(projectID uint) (assessments []model.ProjectAssessment, err error) {
	if err = repository.DB.Where("project_id=?", projectID).
		Preload("Judge").Preload("Criteria").
		Order("judge_id asc, id asc").
		Find(&assessments).Error; err != nil {
		return
	}
	return
}

func (repository *ProjectAssessmentRepositoryImpl) FindByProjectIDAndJudgeID(projectID, judgeID uint) (assessments []model.ProjectAssessment, err error) {
	if err = repository.DB.Where("project_id=? AND judge_id=?", projectID, judgeID).
		Preload("Judge").Preload("Criteria").
		Find(&assessments).Error; err != nil {
		return
	}
	return
}
