package main

import (
	"be-sagara-hackathon/src/cores/database"
	"be-sagara-hackathon/src/middlewares"
	"be-sagara-hackathon/src/modules/auth"
	"be-sagara-hackathon/src/modules/event"
	"be-sagara-hackathon/src/modules/general/upload"
	"be-sagara-hackathon/src/modules/home"
	"be-sagara-hackathon/src/modules/master-data/occupation"
	"be-sagara-hackathon/src/modules/master-data/region"
	"be-sagara-hackathon/src/modules/master-data/skill"
	"be-sagara-hackathon/src/modules/master-data/speciality"
	"be-sagara-hackathon/src/modules/master-data/technology"
	"be-sagara-hackathon/src/modules/payment"
	"be-sagara-hackathon/src/modules/project"
	"be-sagara-hackathon/src/modules/schedule"
	"be-sagara-hackathon/src/modules/team"
	"fmt"

	//"be-sagara-hackathon/src/modules/team"
	"be-sagara-hackathon/src/modules/user"
	"os"

	"github.com/gin-gonic/gin"
)

func InitApp() {
	// Setup Database Connection
	db := database.SetupDatabase()
	//database.MigrateDb(db)
	//seeder.RunSeeder(db)

	// initialize modules/apps
	auth.New(db).InitModule()
	home.New(db).InitModule()
	event.New(db).InitModule()
	speciality.New(db).InitModule()
	occupation.New(db).InitModule()
	skill.New(db).InitModule()
	technology.New(db).InitModule()
	region.New(db).InitModule()
	user.New(db).InitModule()
	payment.New(db).InitModule()
	team.New(db).InitModule()
	schedule.New(db).InitModule()
	project.New(db).InitModule()
	upload.New().InitModule()

	// Get Gin Mode from ENV
	mode := os.Getenv("GIN_MODE")

	// Set Gin Mode
	gin.SetMode(mode)

	// Create New App Instance
	app := gin.Default()

	// Setup CORS
	// app.Use(cors.Default())
	app.Use(middlewares.CORSMiddleware())
	//app.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"https://foo.com"},
	//	AllowMethods:     []string{"PUT", "POST", "GET"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return origin == "https://github.com"
	//	},
	//	MaxAge: 12 * time.Hour,
	//}))

	// Setup Routes
	SetupRoutes(app)

	// Run App at 3000
	err := app.Run(fmt.Sprintf(":%s", os.Getenv("API_PORT")))
	if err != nil {
		return
	}
}
