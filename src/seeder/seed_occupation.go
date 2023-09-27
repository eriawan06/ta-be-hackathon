package seeder

import (
	"be-sagara-hackathon/src/modules/master-data/occupation/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeederOccupation(db *gorm.DB, occupation []model.Occupation) error {
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&occupation).Error; err != nil {
		return err
	}

	return nil
}
