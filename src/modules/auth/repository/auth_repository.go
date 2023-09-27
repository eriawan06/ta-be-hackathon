package repository

import (
	"be-sagara-hackathon/src/modules/auth/model"
	eve "be-sagara-hackathon/src/modules/event/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type AuthRepository interface {
	Register(request model.RegisterModel) (user um.User, err error)
}

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{DB: db}
}

func (repository *AuthRepositoryImpl) Register(request model.RegisterModel) (user um.User, err error) {
	tx := repository.DB.Begin()

	err = tx.Create(&request.User).Error
	var mySqlErr *mysql.MySQLError
	if errors.As(err, &mySqlErr) && mySqlErr.Number == 1062 {
		if strings.Contains(mySqlErr.Message, "idx_unique_user_phone") {
			err = e.ErrPhoneNumberAlreadyExists
		} else if strings.Contains(mySqlErr.Message, "idx_unique_user_email") {
			err = e.ErrEmailAlreadyExists
		}

		tx.Rollback()
		return
	}

	participant := &um.Participant{
		BaseEntity: common.BaseEntity{
			CreatedBy: "self",
			UpdatedBy: "self",
		},
		UserID: request.User.ID,
	}
	if err = tx.Create(&participant).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Create(&eve.EventParticipant{
		BaseEntity: common.BaseEntity{
			CreatedBy: "self",
			UpdatedBy: "self",
		},
		EventID:       request.LatestEventID,
		ParticipantID: participant.ID,
	}).Error; err != nil {
		tx.Rollback()
		return
	}

	request.Invoice.ParticipantID = participant.ID
	if err = tx.Create(&request.Invoice).Error; err != nil {
		tx.Rollback()
		return
	}

	if request.Verification != nil {
		request.Verification.UserID = request.User.ID
		if err = tx.Create(&request.Verification).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	user = request.User
	return
}
