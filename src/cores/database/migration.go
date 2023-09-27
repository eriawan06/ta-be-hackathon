package database

import (
	aum "be-sagara-hackathon/src/modules/auth/model"
	evm "be-sagara-hackathon/src/modules/event/model"
	regm "be-sagara-hackathon/src/modules/master-data/region/model"
	skm "be-sagara-hackathon/src/modules/master-data/skill/model"
	tecm "be-sagara-hackathon/src/modules/master-data/technology/model"
	pym "be-sagara-hackathon/src/modules/payment/model"
	prom "be-sagara-hackathon/src/modules/project/model"
	scm "be-sagara-hackathon/src/modules/schedule/model"
	tm "be-sagara-hackathon/src/modules/team/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"gorm.io/gorm"
)

func MigrateDb(db *gorm.DB) {
	var err error

	err = db.AutoMigrate(&regm.RegProvince{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&regm.RegCity{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&regm.RegDistrict{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&regm.RegVillage{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&skm.Skill{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&tecm.Technology{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&um.UserRole{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&um.User{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&um.Participant{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&um.ParticipantSkill{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&aum.VerificationCode{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&evm.Event{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventMentor{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventJudge{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventCompany{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventTimeline{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventRule{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventFaq{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventAssessmentCriteria{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&evm.EventParticipant{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&pym.Invoice{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&pym.PaymentMethod{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&pym.Payment{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&tm.Team{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&tm.TeamEvent{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&tm.TeamMember{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&tm.TeamInvitation{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&tm.TeamRequest{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&prom.Project{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&prom.ProjectTechnology{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&prom.ProjectSiteLink{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&prom.ProjectImage{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&prom.ProjectAssessment{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&scm.Schedule{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&scm.ScheduleTeam{})
	if err != nil {
		return
	}

	db.Exec("ALTER TABLE specialities ADD CONSTRAINT idx_unique_speciality_name UNIQUE KEY(`name`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE skills ADD CONSTRAINT idx_unique_skill_name UNIQUE KEY(`name`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE occupations ADD CONSTRAINT idx_unique_occupation_name UNIQUE KEY(`name`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE technologies ADD CONSTRAINT idx_unique_technology_name UNIQUE KEY(`name`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE user_roles ADD CONSTRAINT idx_unique_role_name UNIQUE KEY(`name`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE users ADD CONSTRAINT idx_unique_user_email UNIQUE KEY(`email`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE users ADD CONSTRAINT idx_unique_user_phone UNIQUE KEY(`phone_number`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE users ADD CONSTRAINT idx_unique_user_username UNIQUE KEY(`username`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE event_companies ADD CONSTRAINT idx_unique_company_name UNIQUE KEY(`name`, `event_id`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE event_companies ADD CONSTRAINT idx_unique_company_email UNIQUE KEY(`email`, `event_id`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE event_companies ADD CONSTRAINT idx_unique_company_phone UNIQUE KEY(`phone_number`, `event_id`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE teams ADD CONSTRAINT idx_unique_team_code UNIQUE KEY(`code`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
	db.Exec("ALTER TABLE teams ADD CONSTRAINT idx_unique_team_name UNIQUE KEY(`name`, (coalesce(`deleted_at`, '1900-01-01 12:50:18.262000000')));")
}
