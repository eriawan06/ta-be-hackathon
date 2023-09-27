package seeder

import (
	"be-sagara-hackathon/src/modules/master-data/technology/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeederTechnology(db *gorm.DB, technologies []model.Technology) error {
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&technologies).Error; err != nil {
		return err
	}

	return nil
}
