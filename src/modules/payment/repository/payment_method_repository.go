package repository

import (
	"be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type PaymentMethodRepository interface {
	Save(paymentMethod model.PaymentMethod) error
	Update(id uint, paymentMethod model.PaymentMethod) error
	Delete(id uint, deletedBy string) error
	FindAll(
		filter model.FilterPaymentMethod,
		pg *utils.PaginateQueryOffset,
	) (methods []model.PaymentMethod, totalData, totalPage int64, err error)
	FindByID(id uint) (method model.PaymentMethod, err error)
}

type PaymentMethodRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &PaymentMethodRepositoryImpl{DB: db}
}

func (repository *PaymentMethodRepositoryImpl) Save(paymentMethod model.PaymentMethod) error {
	if err := repository.DB.Create(&paymentMethod).Error; err != nil {
		return err
	}
	return nil
}

func (repository *PaymentMethodRepositoryImpl) Update(id uint, paymentMethod model.PaymentMethod) error {
	if err := repository.DB.Select("*").
		Where("id = ?", id).
		Updates(&paymentMethod).Error; err != nil {
		return err
	}
	return nil
}

func (repository *PaymentMethodRepositoryImpl) Delete(paymentMethodId uint, deletedBy string) error {
	if err := repository.DB.Model(&model.PaymentMethod{}).
		Where("id=?", paymentMethodId).
		Updates(map[string]interface{}{"is_active": false, "deleted_at": time.Now(), "deleted_by": deletedBy}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repository *PaymentMethodRepositoryImpl) FindAll(
	filter model.FilterPaymentMethod,
	pg *utils.PaginateQueryOffset,
) (methods []model.PaymentMethod, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&methods).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalPaymentMethod(&filter)
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

func (repository *PaymentMethodRepositoryImpl) getTotalPaymentMethod(filter *model.FilterPaymentMethod) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.PaymentMethod{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterPaymentMethod) (where []string, whereVal []interface{}) {
	where = append(where, "deleted_at IS NULL")

	if filter.Status != "" {
		isActive := false
		if filter.Status == "active" {
			isActive = true
		}
		where = append(where, "is_active = ?")
		whereVal = append(whereVal, isActive)
	}

	return
}

func (repository *PaymentMethodRepositoryImpl) FindByID(id uint) (method model.PaymentMethod, err error) {
	result := repository.DB.
		Where("id=? AND deleted_at IS NULL AND deleted_by IS NULL", id).
		First(&method)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result.Error = e.ErrDataNotFound
	}
	return
}
