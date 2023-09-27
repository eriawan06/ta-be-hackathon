package seeder

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/constants"
	"be-sagara-hackathon/src/utils/helper"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func SeederUser(db *gorm.DB, user model.User) error {
	var checkUser model.User

	db.First(&checkUser, "email=?", user.Email)

	if checkUser.ID == 0 {
		if user.Password != nil {
			hashedPassword, err := utils.HashPassword(helper.DereferString(user.Password))
			if err != nil {
				return errors.New("failed to create user " + user.Email)
			}
			user.Password = helper.ReferString(hashedPassword)
		}

		user.AuthType = constants.AuthTypeRegular
		user.IsActive = true
		user.CreatedBy = "system"
		user.UpdatedBy = "system"
		result := db.Create(&user)
		fmt.Println(result.Error)
		fmt.Println(result.Statement)
		fmt.Println(result.RowsAffected)
		fmt.Println()
		if result.RowsAffected == 0 {
			return errors.New("failed to create user " + user.Email)
		}
	}

	return nil
}
