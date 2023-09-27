package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "be-sagara-hackathon/docs"
	"be-sagara-hackathon/src/middlewares"
	routerAuth "be-sagara-hackathon/src/modules/auth/router"
	routerEvent "be-sagara-hackathon/src/modules/event/router"
	routerUpload "be-sagara-hackathon/src/modules/general/upload/router"
	routerHome "be-sagara-hackathon/src/modules/home/router"
	routerOccupation "be-sagara-hackathon/src/modules/master-data/occupation/router"
	routerRegion "be-sagara-hackathon/src/modules/master-data/region/router"
	routerSkill "be-sagara-hackathon/src/modules/master-data/skill/router"
	routerSpeciality "be-sagara-hackathon/src/modules/master-data/speciality/router"
	routerTechnology "be-sagara-hackathon/src/modules/master-data/technology/router"
	routerPayment "be-sagara-hackathon/src/modules/payment/router"
	routerProject "be-sagara-hackathon/src/modules/project/router"
	routerSchedule "be-sagara-hackathon/src/modules/schedule/router"
	routerTeam "be-sagara-hackathon/src/modules/team/router"
	routerUser "be-sagara-hackathon/src/modules/user/router"
	//teamRouter "be-sagara-hackathon/src/modules/team/router"
)

// SetupRoutes Setup Routes
func SetupRoutes(app *gin.Engine) {
	// Check Server Status Endpoint
	app.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Server alive!",
			"data":    context,
		})
	})

	// Use ginSwagger middleware to serve the API docs
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	// Setup Routes Group
	auth := app.Group("/api/v1/auth")
	{
		routerAuth.RegisterRoutes(auth)
	}

	home := app.Group("/api/v1/home")
	{
		routerHome.HomeRouter(home)
	}

	region := app.Group("/api/v1/region")
	{
		routerRegion.RegionRouter(region)
	}

	v1 := app.Group("/api/v1")
	{
		v1.Use(middlewares.JwtAuthMiddleware())

		routerSpeciality.SpecialityRouter(v1.Group("/specialities"))
		routerOccupation.OccupationRouter(v1.Group("/occupations"))
		routerSkill.SkillRouter(v1.Group("/skills"))
		routerTechnology.TechnologyRouter(v1.Group("/technologies"))
		routerUser.UserRouter(v1.Group("/users"))
		routerEvent.EventRouter(v1.Group("/events"))
		routerPayment.PaymentMethodRouter(v1.Group("/payment-methods"))
		routerPayment.InvoiceRouter(v1.Group("/invoices"))
		routerPayment.PaymentRouter(v1.Group("/payments"))
		routerTeam.TeamRouter(v1.Group("/teams"))
		routerProject.ProjectRouter(v1.Group("/projects"))
		routerSchedule.ScheduleRouter(v1.Group("/schedules"))
		routerUpload.UploadRouter(v1.Group("/upload"))
	}
}
