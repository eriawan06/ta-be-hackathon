package repository

import (
	aum "be-sagara-hackathon/src/modules/auth/model"
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user model.User) error
	Update(userID uint, user model.User) error
	UpdateByEmail(user model.User) error
	UpdateStatus(userId uint, isActive bool) error
	UpdatePassword(userId uint, password, verifCode string) error
	Delete(userID uint, deletedBy string) error
	Find(
		filter model.FilterUser,
		pg *utils.PaginateQueryOffset,
	) (users []model.User, totalData, totalPage int64, err error)
	FindByID(id uint) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByUsername(username string) (model.User, error)
	FindUserProfileByEmail(email string) (model.UserProfile, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (repository *UserRepositoryImpl) Save(user model.User) error {
	result := repository.DB.Create(&user)

	var mySqlErr *mysql.MySQLError
	if errors.As(result.Error, &mySqlErr) && mySqlErr.Number == 1062 {
		if strings.Contains(mySqlErr.Message, "idx_unique_user_phone") {
			result.Error = e.ErrPhoneNumberAlreadyExists
		} else if strings.Contains(mySqlErr.Message, "idx_unique_user_email") {
			result.Error = e.ErrEmailAlreadyExists
		} else if strings.Contains(mySqlErr.Message, "idx_unique_user_username") {
			result.Error = e.ErrUsernameAlreadyExists
		}
	}

	return result.Error
}

func (repository *UserRepositoryImpl) Update(userID uint, user model.User) error {
	if err := repository.DB.Model(&model.User{}).Select("*").
		Where("id = ?", userID).
		Updates(&user).
		Error; err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) && mySqlErr.Number == 1062 {
			if strings.Contains(mySqlErr.Message, "idx_unique_user_phone") {
				err = e.ErrPhoneNumberAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_user_email") {
				err = e.ErrEmailAlreadyExists
			} else if strings.Contains(mySqlErr.Message, "idx_unique_user_username") {
				err = e.ErrUsernameAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (repository *UserRepositoryImpl) UpdateByEmail(user model.User) error {
	result := repository.DB.Model(&model.User{}).Where("email = ?", &user.Email).Updates(&user)
	return result.Error
}

func (repository *UserRepositoryImpl) UpdateStatus(userId uint, isActive bool) error {
	result := repository.DB.Model(&model.User{}).Where("id=?", userId).Update("is_active", isActive)
	return result.Error
}

func (repository *UserRepositoryImpl) UpdatePassword(userId uint, password, verifCode string) error {
	tx := repository.DB.Begin()
	if err := tx.Model(&model.User{}).
		Where("id=?", userId).
		Update("password", password).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(verifCode) > 0 {
		if err := tx.Where("code=?", verifCode).
			Delete(&aum.VerificationCode{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (repository *UserRepositoryImpl) Delete(userID uint, deletedBy string) error {
	if err := repository.DB.Model(&model.User{}).
		Where("id=?", userID).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "deleted_by": deletedBy}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositoryImpl) Find(
	filter model.FilterUser,
	pg *utils.PaginateQueryOffset,
) (users []model.User, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilter(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.
		Preload("Participant").
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&users).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalUser(&filter)
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

func (repository *UserRepositoryImpl) getTotalUser(filter *model.FilterUser) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilter(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Model(&model.User{}).
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilter(filter model.FilterUser) (where []string, whereVal []interface{}) {
	where = append(where, "deleted_at is null")

	if filter.ExceptUserID != 0 {
		where = append(where, "id != @except_user_id")
		whereVal = append(whereVal, sql.Named("except_user_id", filter.ExceptUserID))
	}

	if filter.RoleID != 0 {
		where = append(where, "user_role_id = @role")
		whereVal = append(whereVal, sql.Named("role", filter.RoleID))
	}

	if filter.Status != "" {
		isActive := false
		if filter.Status == "active" {
			isActive = true
		}
		where = append(where, "is_active = @is_active")
		whereVal = append(whereVal, sql.Named("is_active", isActive))
	}

	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		where = append(where, "(LOWER(name) LIKE @keyword OR LOWER(email) LIKE @keyword OR phone_number LIKE @keyword )")
		whereVal = append(whereVal, sql.Named("keyword", "%"+filter.Search+"%"))
	}

	return
}

func (repository *UserRepositoryImpl) FindByID(id uint) (model.User, error) {
	var user model.User

	if err := repository.DB.
		Where("id = ?", id).
		Preload("UserRole").
		Preload("Occupation").
		Preload("Participant").
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return user, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindByEmail(email string) (model.User, error) {
	var user model.User

	result := repository.DB.
		Where("email = ? AND deleted_by is null", email).
		Preload("UserRole").
		Preload("Participant").
		Preload("Participant.Skills").
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result.Error = e.ErrEmailNotRegistered
		}
		return user, result.Error
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindByUsername(username string) (model.User, error) {
	var user model.User

	result := repository.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindUserProfileByEmail(email string) (model.UserProfile, error) {
	query := `SELECT
		u.name as full_name, u.email, u.phone_number, u.institution, 
		ur.id as role_id, ur.name as role_name, u.is_active
		FROM users u
		LEFT JOIN user_roles ur ON ur.id = u.user_role_id
		WHERE u.email = ?
		LIMIT 1`

	var userProfile model.UserProfile
	if err := repository.DB.Raw(query, email).Scan(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return userProfile, err
	}

	return userProfile, nil
}
