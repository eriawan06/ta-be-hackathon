package repository

import (
	"be-sagara-hackathon/src/modules/event/model"
	e "be-sagara-hackathon/src/utils/errors"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type EventCompanyRepository interface {
	Save(ec model.EventCompany) error
	Update(ecID uint, ec model.EventCompany) error
	Delete(ecID uint) error
	FindAll(filter model.FilterEventCompany) (companies []model.EventCompany, err error)
	FindOne(ecID uint) (company model.EventCompany, err error)
	FindManyByEventID(eventID uint) (companies []model.EventCompany, err error)
}

type EventCompanyRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventCompanyRepository(db *gorm.DB) EventCompanyRepository {
	return &EventCompanyRepositoryImpl{DB: db}
}

func (repository *EventCompanyRepositoryImpl) Save(ec model.EventCompany) error {
	if err := repository.DB.Create(&ec).Error; err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) && mySqlErr.Number == 1062 {
			if strings.Contains(mySqlErr.Message, "idx_unique_company_phone") {
				err = e.ErrPhoneNumberAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_company_email") {
				err = e.ErrEmailAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_company_name") {
				err = e.ErrNameAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (repository *EventCompanyRepositoryImpl) Update(ecID uint, ec model.EventCompany) error {
	if err := repository.DB.Select("*").Where("id = ?", ecID).Updates(&ec).Error; err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) && mySqlErr.Number == 1062 {
			if strings.Contains(mySqlErr.Message, "idx_unique_company_phone") {
				err = e.ErrPhoneNumberAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_company_email") {
				err = e.ErrEmailAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_company_name") {
				err = e.ErrNameAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (repository *EventCompanyRepositoryImpl) Delete(ecID uint) error {
	if err := repository.DB.Delete(&model.EventCompany{}, ecID).Error; err != nil {
		return err
	}
	return nil
}

func (repository *EventCompanyRepositoryImpl) FindAll(filter model.FilterEventCompany) (companies []model.EventCompany, err error) {
	if err = repository.DB.Where("event_id=?", filter.EventID).Find(&companies).Error; err != nil {
		return
	}
	return
}

func (repository *EventCompanyRepositoryImpl) FindOne(ecID uint) (company model.EventCompany, err error) {
	if err = repository.DB.Where("id=?", ecID).First(&company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *EventCompanyRepositoryImpl) FindManyByEventID(eventID uint) (companies []model.EventCompany, err error) {
	if err = repository.DB.Where("event_id=?", eventID).Find(&companies).Error; err != nil {
		return
	}
	return
}
