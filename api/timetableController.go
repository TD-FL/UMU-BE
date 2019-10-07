package api

import (
	"fuu-be/visualizer"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// TimetableByCourse godoc
// @Summary Returns the schedule for a specific course.
// @Description TimetableByCourse
// @Produce  json
// @Param faculty path string true "Faculty"
// @Param name path string true "Name of course"
// @Success 200 {array} visualizer.Event
// @Router /schedule/course/{faculty}/{name} [get]
func TimetableByCourse(c *gin.Context) {
	name := c.Param("name")
	faculty := c.Param("faculty")

	aEvent := make([]visualizer.Event, 0)
	aCourses := strings.Split(name, ",")

	for _, v := range aCourses {
		aEvent = append(aEvent, visualizer.FacultyTimeTables[faculty].EventsByCourses[v]...)
	}

	result := aEvent

	if result != nil && len(result) > 0 {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// TimetableByProfessor godoc
// @Summary Returns the schedule for a specific professor.
// @Description TimetableByProfessor
// @Produce  json
// @Param faculty path string true "Faculty"
// @Param name path string true "Name of professor"
// @Success 200 {array} visualizer.Event
// @Router /schedule/professor/{faculty}/{type} [get]
func TimetableByProfessor(c *gin.Context) {
	name := c.Param("name")
	faculty := c.Param("faculty")

	result := visualizer.FacultyTimeTables[faculty].EventsByProfessors[name]

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// TimetableByGroupType godoc
// @Summary Returns the schedule for a specific Group Type.
// @Description TimetableByGroupType
// @Produce  json
// @Param faculty path string true "Faculty"
// @Param type path string true "Type"
// @Success 200 {array} visualizer.Event
// @Router /schedule/faculty/{faculty}/{type} [get]
func TimetableByGroupType(c *gin.Context) {
	gType := c.Param("type")
	faculty := c.Param("faculty")

	result := visualizer.FacultyTimeTables[faculty].EventsByGroups[gType]

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// TimetableByGroupTypeAndYear godoc
// @Summary Returns the schedule for a specific Group Type and year.
// @Description TimetableByGroupTypeAndYear
// @Produce  json
// @Param faculty path string true "Faculty"
// @Param type path string true "Type"
// @Param year path string true "Year"
// @Success 200 {array} visualizer.Event
// @Router /schedule/faculty/{faculty}/{type}/{year} [get]
func TimetableByGroupTypeAndYear(c *gin.Context) {
	gType := c.Param("type")
	faculty := c.Param("faculty")
	year, err := strconv.ParseInt(c.Param("year"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad request!")
	}

	result := visualizer.FacultyTimeTables[faculty].EventsByGroups[gType][year]

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// GetAllOfFaculty godoc
// @Summary Returns the schedules of all groups in the faculty.
// @Description GetAllOfFaculty
// @Produce  json
// @Param faculty path string true "Faculty"
// @Success 200  {array} visualizer.Event
// @Router /schedule/faculty/{faculty} [get]
func GetAll(c *gin.Context) {
	faculty := c.Param("faculty")

	result := visualizer.FacultyTimeTables[faculty].EventsByGroups

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}
