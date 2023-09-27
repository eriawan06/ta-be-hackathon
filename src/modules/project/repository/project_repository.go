package repository

import (
	"be-sagara-hackathon/src/modules/project/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"strings"
)

type ProjectRepository interface {
	Create(project model.Project) error
	Update(id uint, req model.UpdateProjectModel) error
	UpdateStatus(id uint, project model.Project) error
	FindOne(id uint) (project model.Project, err error)
	FindAll(
		filter model.FilterProject,
		pg *utils.PaginateQueryOffset,
	) (projects []model.ProjectLite, totalData, totalPage int64, err error)
	FindByEventIDAndTeamID(eventID, teamID uint) (project model.Project, err error)
}

type ProjectRepositoryImpl struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &ProjectRepositoryImpl{DB: db}
}

func (repository *ProjectRepositoryImpl) Create(project model.Project) error {
	tx := repository.DB.Begin()
	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repository *ProjectRepositoryImpl) Update(id uint, req model.UpdateProjectModel) error {
	tx := repository.DB.Begin()
	if err := tx.Select("*").
		Omit("BuiltWith", "SiteLinks", "Images").
		Where("id=?", id).
		Updates(&req.Project).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(req.Project.BuiltWith) > 0 {
		if err := tx.Create(&req.Project.BuiltWith).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(req.Project.SiteLinks) > 0 {
		if err := tx.Create(&req.Project.SiteLinks).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(req.Project.Images) > 0 {
		if err := tx.Create(&req.Project.Images).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(req.RemovedBuiltWith) > 0 {
		if err := tx.Delete(&req.RemovedBuiltWith).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(req.RemovedSiteLinks) > 0 {
		if err := tx.Delete(&model.ProjectSiteLink{}, req.RemovedSiteLinks).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(req.RemovedImages) > 0 {
		if err := tx.Delete(&model.ProjectImage{}, req.RemovedImages).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (repository *ProjectRepositoryImpl) UpdateStatus(id uint, project model.Project) error {
	if err := repository.DB.Table("projects").Where("id=?", id).
		Updates(map[string]interface{}{
			"status":     project.Status,
			"updated_at": project.UpdatedAt,
			"updated_by": project.UpdatedBy,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (repository *ProjectRepositoryImpl) FindOne(id uint) (project model.Project, err error) {
	if err = repository.DB.Preload(clause.Associations).
		Preload("BuiltWith.Technology").
		Where("id=?", id).
		First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *ProjectRepositoryImpl) FindAll(
	filter model.FilterProject,
	pg *utils.PaginateQueryOffset,
) (projects []model.ProjectLite, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("projects").
		Select(`id, event_id, team_id, name, status, created_at`).
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&projects).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalProject(&filter)
	if err != nil {
		return
	}

	if pg.Limit > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(pg.Limit)))
	} else {
		totalPage = 1
	}

	return
}

func (repository *ProjectRepositoryImpl) getTotalProject(filter *model.FilterProject) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.Project{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterProject) (where []string, whereVal []interface{}) {
	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		where = append(where, "LOWER(name) LIKE ?")
		whereVal = append(whereVal, "%"+filter.Search+"%")
	}

	if filter.Status != "" {
		where = append(where, "status = ?")
		whereVal = append(whereVal, filter.Status)
	}

	if filter.CreatedAt != "" {
		where = append(where, "date(created_at) = ?")
		whereVal = append(whereVal, filter.CreatedAt)
	}

	if filter.TeamID != 0 {
		where = append(where, "team_id = ?")
		whereVal = append(whereVal, filter.TeamID)
	}

	if filter.EventID != 0 {
		where = append(where, "event_id = ?")
		whereVal = append(whereVal, filter.EventID)
	}

	return
}

func (repository *ProjectRepositoryImpl) FindByEventIDAndTeamID(eventID, teamID uint) (project model.Project, err error) {
	if err = repository.DB.Where("event_id=? AND team_id=?", eventID, teamID).
		First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
