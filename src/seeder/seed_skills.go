package seeder

import (
	"be-sagara-hackathon/src/modules/master-data/skill/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeederSkills(db *gorm.DB, skill []model.Skill) error {
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&skill).Error; err != nil {
		return err
	}

	return nil
}
