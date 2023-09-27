package seeder

import (
	"be-sagara-hackathon/src/modules/user/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func SeederRole(db *gorm.DB, role model.UserRole) error {
	var checkRole model.UserRole

	db.First(&checkRole, "name=?", role.Name)

	if checkRole.ID == 0 {
		// result := db.Exec("INSERT INTO user_roles (name) VALUES (?)", name)
		role.CreatedBy = "system"
		role.UpdatedBy = "system"
		result := db.Create(&role)
		fmt.Println(result.Error)
		fmt.Println(result.Statement)
		fmt.Println(result.RowsAffected)
		fmt.Println()
		if result.RowsAffected == 0 {
			return errors.New("failed to create role " + role.Name)
		}
	}

	return nil
}
