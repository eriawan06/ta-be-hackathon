package repository

import (
	"be-sagara-hackathon/src/modules/user/model"
	e "be-sagara-hackathon/src/utils/errors"
	"gorm.io/gorm"
)

type UserRoleRepository interface {
	FindByName(name string) (role model.UserRole, err error)
}

type UserRoleRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &UserRoleRepositoryImpl{DB: db}
}

func (repository *UserRoleRepositoryImpl) FindByName(name string) (role model.UserRole, err error) {
	if err = repository.DB.Where("name=?", name).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}
