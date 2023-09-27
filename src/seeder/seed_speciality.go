package seeder

import (
	"be-sagara-hackathon/src/modules/master-data/speciality/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeederSpeciality(db *gorm.DB, speciality []model.Speciality) error {
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&speciality).Error; err != nil {
		return err
	}

	return nil
}
