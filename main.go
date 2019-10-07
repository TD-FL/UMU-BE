package main

import (
	"fmt"
	"fuu-be/api"
	"fuu-be/downloader"
	"fuu-be/parser"
	"fuu-be/visualizer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"io/ioutil"
	"time"

	_ "fuu-be/docs"
)

// @title Swagger za UMU REST API
// @version 2.0
// @description Urnik Mariborske Univerze (UMU) je projekt katerega cilj je, da študentom UM nudi boljšo storitev za urnik kot je WiseTimeTables.
// @termsOfService http://swagger.io/terms/

// @license.name GPL 2.0
// @license.url https://opensource.org/licenses/gpl-2.0.php

// @host api.urnik-mb.cf
// @BasePath /api/v2

// @x-extension-openapi {"example": "value on a json format"}

const API2_URI = "/api/v2/"
const privacyPolicy_URI = "/privacy"

const timetableURL = "http://www.wisetimetable.com/m/"

func main() {
	//gin.SetMode(gin.ReleaseMode)

	privacy, err := ioutil.ReadFile("privacy_policy.md")
	if err != nil {
		fmt.Print(err)
	}

	api.PrivacyText = string(privacy)

	parser.ParseFaculty()
	schedule(getData, 2*time.Hour)

	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(cors.Default())

	//privacy location
	router.GET(privacyPolicy_URI, api.GetPrivacy)

	//swagger
	url := ginSwagger.URL("https://api.urnik-mb.cf/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET(API2_URI+"week", api.GetCurrentWeek)
	router.GET(API2_URI+"faculties", api.GetFaculties)

	router.GET(API2_URI+"groups/:faculty", api.GetGroups)
	router.GET(API2_URI+"groups/:faculty/:years", api.GetGroupsWithYears)
	router.GET(API2_URI+"courses/:faculty", api.GetKnownCourses)
	router.GET(API2_URI+"professors/:faculty", api.GetKnownProfessors)
	router.GET(API2_URI+"years/:faculty/:type/", api.GetYearsOfType)

	router.GET(API2_URI+"schedule/course/:faculty/:name", api.TimetableByCourse)
	router.GET(API2_URI+"schedule/professor/:faculty/:name", api.TimetableByProfessor)

	router.GET(API2_URI+"schedule/faculty/:faculty", api.GetAll)
	router.GET(API2_URI+"schedule/faculty/:faculty/:type", api.TimetableByGroupType)
	router.GET(API2_URI+"schedule/faculty/:faculty/:type/:year", api.TimetableByGroupTypeAndYear)

	router.Run(":8089")
}

func clearData() {
	parser.Clean()
	visualizer.Clean()
}

func getData() {
	clearData()

	for _, v := range parser.FacultyMap {
		err := downloader.Download("./timetable.zip", timetableURL+v.Hash+".zip")
		if err != nil {
			panic(err)
		}
		downloader.Extract("./timetable.zip")

		parser.Parse()

		visualizer.Vizualize(v.ShortName)
	}

	if gin.Mode() == gin.DebugMode {
		fmt.Println("[GIN-debug] Parsing done!")
	}
}

func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
