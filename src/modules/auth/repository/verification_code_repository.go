package repository

import (
	"be-sagara-hackathon/src/modules/auth/model"
	e "be-sagara-hackathon/src/utils/errors"
	"errors"
	"gorm.io/gorm"
)

type VerificationCodeRepository interface {
	Save(vc model.VerificationCode) error
	Delete(vcID, userID uint) error
	DeleteByCode(code string) error
	FindOne(vcId uint) (model.VerificationCode, error)
	FindByCode(code string) (model.VerificationCode, error)
	FindByUserIdAndCodeType(userId uint, codeType string) (model.VerificationCode, error)
	FindByUserEmailAndCodeType(email, codeType string) (model.VerificationCode, error)
}

type VerificationCodeRepositoryImpl struct {
	DB *gorm.DB
}

func NewVerificationCodeRepository(db *gorm.DB) VerificationCodeRepository {
	return &VerificationCodeRepositoryImpl{DB: db}
}

func (repository *VerificationCodeRepositoryImpl) Save(vc model.VerificationCode) error {
	result := repository.DB.Create(&vc)
	return result.Error
}

func (repository *VerificationCodeRepositoryImpl) Delete(vcID, userID uint) error {
	tx := repository.DB.Begin()
	if err := tx.Delete(&model.VerificationCode{}, vcID).Error; err != nil {
		tx.Rollback()
		return err
	}

	//update user's status to 'active'
	if userID != 0 {
		if err := tx.Table("users").
			Where("id=?", userID).
			Update("is_active", true).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (repository *VerificationCodeRepositoryImpl) DeleteByCode(code string) error {
	result := repository.DB.Where("code=?", code).Delete(&model.VerificationCode{})
	return result.Error
}

func (repository *VerificationCodeRepositoryImpl) FindOne(vcId uint) (model.VerificationCode, error) {
	var verificationCode model.VerificationCode
	result := repository.DB.Where("id=?", vcId).First(&verificationCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result.Error = e.ErrDataNotFound
	}
	return verificationCode, result.Error
}

func (repository *VerificationCodeRepositoryImpl) FindByCode(code string) (model.VerificationCode, error) {
	var verificationCode model.VerificationCode
	result := repository.DB.Where("code=?", code).First(&verificationCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result.Error = e.ErrDataNotFound
	}
	return verificationCode, result.Error
}

func (repository *VerificationCodeRepositoryImpl) FindByUserIdAndCodeType(userId uint, codeType string) (model.VerificationCode, error) {
	var verificationCode model.VerificationCode
	result := repository.DB.Where("user_id=? AND type=?", userId, codeType).First(&verificationCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result.Error = e.ErrDataNotFound
	}
	return verificationCode, result.Error
}

func (repository *VerificationCodeRepositoryImpl) FindByUserEmailAndCodeType(email, codeType string) (model.VerificationCode, error) {
	var verificationCode model.VerificationCode
	query := `
		SELECT vc.*, u.email
		FROM verification_codes vc
		LEFT JOIN users u ON u.id = vc.user_id
		WHERE u.email=? AND vc.type=?
	`
	result := repository.DB.Raw(query, email, codeType).Scan(&verificationCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result.Error = e.ErrDataNotFound
	}
	return verificationCode, result.Error
}
