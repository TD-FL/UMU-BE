package api

import (
	"fuu-be/parser"
	"fuu-be/visualizer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCurrentWeek godoc
// @Summary Returns the current week number.
// @Description GetCurrentWeek
// @Produce  json
// @Success 200 {object} visualizer.CurrentWeek
// @Router /week [get]
func GetCurrentWeek(c *gin.Context) {
	c.JSON(http.StatusOK, visualizer.GetCurrentWeek())
}

// GetFaculties godoc
// @Summary Returns the list of all available faculties.
// @ID get-string-by-int
// @Description GetFaculties
// @Produce  json
// @Success 200 {array} parser.Faculty
// @Router /faculties [get]
func GetFaculties(c *gin.Context) {
	result := parser.GetKnownFaculties()

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// GetKnownCourses godoc
// @Summary Returns the list of all available courses in a specific faculty.
// @Description GetKnownCourses
// @Produce  json
// @Param faculty path string true "Faculty"
// @Success 200 {array} string
// @Router /courses/{faculty} [get]
func GetKnownCourses(c *gin.Context) {
	faculty := c.Param("faculty")

	result := visualizer.GetKnownCourses(faculty)

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// GetYearsOfType godoc
// @Summary Returns an array of years in a specific group type.
// @Description GetYearsOfType
// @Produce  json
// @Param faculty path string true "Faculty"
// @Param type path string true "Type"
// @Success 200 {array} int64
// @Router /years/{faculty}/{type} [get]
func GetYearsOfType(c *gin.Context) {
	gType := c.Param("type")
	faculty := c.Param("faculty")

	var keys = make([]int64, 0)

	for k, _ := range visualizer.FacultyTimeTables[faculty].EventsByGroups[gType] {
		keys = append(keys, k)
	}

	result := keys

	if result != nil && len(result) != 0 {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// GetKnownProfessors godoc
// @Summary Returns an array of professors in a specific faculty.
// @Description GetKnownProfessors
// @Produce  json
// @Param faculty path string true "Faculty"
// @Success 200 {array} string
// @Router /professors/{faculty} [get]
func GetKnownProfessors(c *gin.Context) {
	faculty := c.Param("faculty")

	result := visualizer.GetKnownProfessors(faculty)

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// GetGroups godoc
// @Summary Returns an array of Groups in a specific faculty.
// @Description GetGroups
// @Produce  json
// @Param faculty path string true "Faculty"
// @Success 200 {array} string
// @Router /groups/{faculty} [get]
func GetGroups(c *gin.Context) {
	faculty := c.Param("faculty")

	result := visualizer.GetKnownGroups(faculty)

	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}

// GetGroupsWithYears godoc
// @Summary Returns an array of Groups with years in a specific faculty.
// @Description GetGroupsWithYears
// @Produce  json
// @Param faculty path string true "Faculty"
// @Success 200 {array} visualizer.GroupWithYears
// @Router /groups/{faculty}/years [get]
func GetGroupsWithYears(c *gin.Context) {
	faculty := c.Param("faculty")
	years := c.Param("years")

	result := visualizer.GetKnownGroupsWithYears(faculty)

	if result != nil && years == "years" { // a necessary evil of GIN
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, "No data!")
	}
}
