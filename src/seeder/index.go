package seeder

import (
	ocm "be-sagara-hackathon/src/modules/master-data/occupation/model"
	skm "be-sagara-hackathon/src/modules/master-data/skill/model"
	spm "be-sagara-hackathon/src/modules/master-data/speciality/model"
	tecm "be-sagara-hackathon/src/modules/master-data/technology/model"
	ue "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/constants"
	"be-sagara-hackathon/src/utils/helper"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func RunSeeder(db *gorm.DB) {
	for _, seed := range SeedsAll() {
		if err := seed.Run(db); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}

// SeedsAll will call all your filled func seeds
func SeedsAll() []Seed {

	return []Seed{
		{
			Name: "SeederRole",
			Run: func(db *gorm.DB) error {
				roles := []ue.UserRole{
					{Name: constants.UserSuperadmin},
					{Name: constants.UserAdmin},
					{Name: constants.UserHR},
					{Name: constants.UserParticipant},
					{Name: constants.UserCompany},
					{Name: constants.UserMentor},
					{Name: constants.UserJudge},
				}
				for _, role := range roles {
					err := SeederRole(db, role)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Name: "SeederUser",
			Run: func(db *gorm.DB) error {
				users := []ue.User{
					{
						Name:        "Superadmin",
						Email:       "sagarahackathon2021@gmail.com",
						Institution: helper.ReferString("Sagara Technology"),
						Password:    helper.ReferString(os.Getenv("SUPERADMIN_PWD")),
						UserRoleID:  1,
					},
				}
				for _, user := range users {
					err := SeederUser(db, user)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Name: "SeederSpeciality",
			Run: func(db *gorm.DB) error {
				jsonData, err := ParseJSONFile("specialities.json")
				if err != nil {
					return err
				}

				var specialities []spm.Speciality
				for _, v := range jsonData.([]interface{}) {
					occupation := spm.Speciality{Name: v.(string)}
					occupation.CreatedBy, occupation.UpdatedBy = "system", "system"
					specialities = append(specialities, occupation)
				}

				_ = SeederSpeciality(db, specialities)
				return nil
			},
		},
		{
			Name: "SeederOccupation",
			Run: func(db *gorm.DB) error {
				jsonData, err := ParseJSONFile("occupations.json")
				if err != nil {
					return err
				}

				var occupations []ocm.Occupation
				for _, v := range jsonData.([]interface{}) {
					occupation := ocm.Occupation{Name: strings.ToTitle(v.(string))}
					occupation.CreatedBy, occupation.UpdatedBy = "system", "system"
					occupations = append(occupations, occupation)
				}

				_ = SeederOccupation(db, occupations)
				return nil
			},
		},
		{
			Name: "SeederSkill",
			Run: func(db *gorm.DB) error {
				jsonData, err := ParseJSONFile("skills.json")
				if err != nil {
					return err
				}

				var skills []skm.Skill
				for _, v := range jsonData.([]interface{}) {
					skill := skm.Skill{Name: strings.ToTitle(v.(string))}
					skill.CreatedBy, skill.UpdatedBy = "system", "system"
					skills = append(skills, skill)
				}

				_ = SeederSkills(db, skills)
				return nil
			},
		},
		{
			Name: "SeederTechnology",
			Run: func(db *gorm.DB) error {
				_, b, _, _ := runtime.Caller(0)
				basePath := filepath.Join(filepath.Dir(b), "../..")
				path := filepath.Join(basePath, "seed-data/technologies")
				files, err := ioutil.ReadDir(path)
				if err != nil {
					log.Fatal(err)
				}

				var filenames []string
				for _, file := range files {
					filenames = append(filenames, file.Name())
				}

				var technologies []tecm.Technology
				for i, filename := range filenames {
					jsonData, err2 := ParseJSONFile(filepath.Join("technologies", filename))
					if err2 != nil {
						return err2
					}

					for k := range jsonData.(map[string]interface{}) {
						technology := tecm.Technology{Name: k}
						technology.CreatedBy, technology.UpdatedBy = "system", "system"
						technologies = append(technologies, technology)
					}

					if (i > 0 && i%4 == 0) || i == len(filenames)-1 {
						_ = SeederTechnology(db, technologies)
						technologies = nil
					}
				}

				return nil
			},
		},
	}
}

func ParseJSONFile(filename string) (output interface{}, err error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")
	path := filepath.Join(basePath, "seed-data", filename)
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, &output); err != nil {
		return
	}
	return
}
