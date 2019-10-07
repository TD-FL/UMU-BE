package visualizer

import (
	"fmt"
	"fuu-be/parser"
	"github.com/gin-gonic/gin"
	"time"
)

type Group struct {
	Field    string
	Year     int64
	Type     string
	SubGroup string
}

type Event struct {
	DayOfWeek int64
	BeginWeek int64
	EndWeek   int64
	StartTime string
	Duration  int64
	Type      string
	Course    string
	Room      string
	Professor string
	Group     Group
}

type CurrentWeek struct {
	CurrentWeek int64
}

type GroupWithYears struct {
	Name  string
	Years []int64
}

type FacultyTimeTable struct {
	Events             []Event
	EventsByGroups     map[string]map[int64][]Event
	EventsByCourses    map[string][]Event
	EventsByProfessors map[string][]Event

	KnownGroups     map[string]string
	KnownCourses    map[string]string
	KnownProfessors map[string]string
}

var FacultyTimeTables = make(map[string]FacultyTimeTable)

var Events = make([]Event, 0)

var EventsByGroups = make(map[string]map[int64][]Event)
var EventsByCourses = make(map[string][]Event)
var EventsByProfessors = make(map[string][]Event)

var KnownGroups = make(map[string]string)
var KnownCourses = make(map[string]string)
var KnownProfessors = make(map[string]string)

func Vizualize(faculty string) {
	if gin.Mode() == gin.DebugMode {
		fmt.Println("[GIN-debug] Parsing faculty '" + faculty + "'")
	}

	if len(parser.Schedules) == 0 {
		return
	}

	for _, schedule := range parser.Schedules {

		if len(schedule.GroupIds) > 0 {
			var groups []Group
			for i, _ := range schedule.GroupIds {
				var group Group

				group = groupGroups(parser.GroupsMap[schedule.GroupIds[i]])

				groups = append(groups, group)
			}

			for _, group := range groups {
				var event Event

				event.DayOfWeek = schedule.DayOfWeek
				event.BeginWeek = schedule.BeginWeek
				event.EndWeek = schedule.EndWeek
				event.StartTime = parser.TimeSlotsMap[schedule.BeginTimeSlotId].StartTimeString
				event.Duration = schedule.Duration
				event.Type = parser.ExecTypesMap[schedule.ExecTypeId].ExecType
				event.Course = parser.CoursesMap[schedule.CourseId].Name
				event.Room = parser.RoomsMap[schedule.RoomId].RoomName

				KnownCourses[event.Course] = event.Course

				if len(schedule.TutorIds) > 0 {
					event.Professor = parser.TutorsMap[schedule.TutorIds[0]].Firstname + " " + parser.TutorsMap[schedule.TutorIds[0]].Lastname
				}

				KnownProfessors[event.Professor] = event.Professor

				event.Group = group

				Events = append(Events, event)

				sortByGroups(event)
				sortByCourse(event)
				sortByProfessor(event)
			}
		}
	}

	facutyTT := FacultyTimeTable{Events, EventsByGroups, EventsByCourses,
		EventsByProfessors, KnownGroups, KnownCourses, KnownProfessors}

	FacultyTimeTables[faculty] = facutyTT

	parser.Clean()
	Clean()
}

func groupGroups(pGroup parser.Group) Group {
	var group Group

	branch := parser.BranchesMap[pGroup.BranchId]
	group.Field = parser.ProgramsMap[branch.ProgramId].Name

	group.Year = branch.Year

	if pGroup.ParentGroupId == -1 {
		group.Type = pGroup.Name
	} else {
		group.Type = parser.GroupsMap[pGroup.ParentGroupId].Name
		group.SubGroup = pGroup.Name
	}

	if group.SubGroup == "" {
		group.SubGroup = group.Type
	}

	return group
}

func Clean() {
	Events = make([]Event, 0)

	EventsByGroups = make(map[string]map[int64][]Event)
	EventsByCourses = make(map[string][]Event)
	EventsByProfessors = make(map[string][]Event)

	KnownGroups = make(map[string]string)
	KnownCourses = make(map[string]string)
	KnownProfessors = make(map[string]string)
}

func sortByGroups(e Event) {
	if e.Room == "" {
		return
	}

	name := e.Group.Field

	if !existsInKnownGroups(name) {
		KnownGroups[name] = name

		EventsByGroups[name] = make(map[int64][]Event)
	}

	EventsByGroups[name][e.Group.Year] = append(EventsByGroups[name][e.Group.Year], e)
}

func sortByCourse(e Event) {
	if e.Room == "" {
		return
	}

	name := e.Course

	EventsByCourses[name] = append(EventsByCourses[name], e)
}

func sortByProfessor(e Event) {
	if e.Room == "" {
		return
	}

	name := e.Professor

	EventsByProfessors[name] = append(EventsByProfessors[name], e)
}

func existsInKnownGroups(s string) bool {
	return KnownGroups[s] != ""
}

func GetKnownGroups(faculty string) []string {
	var groups []string
	for _, v := range FacultyTimeTables[faculty].KnownGroups {
		groups = append(groups, v)
	}

	return groups
}

func GetKnownGroupsWithYears(faculty string) []GroupWithYears {
	var groups []GroupWithYears

	for k, v := range FacultyTimeTables[faculty].EventsByGroups {
		var gwy GroupWithYears
		gwy.Name = k

		gwy.Years = make([]int64, 0)
		for k := range v {
			gwy.Years = append(gwy.Years, k)
		}

		groups = append(groups, gwy)
	}

	return groups
}

func GetKnownCourses(faculty string) []string {
	var courses []string
	for _, v := range FacultyTimeTables[faculty].KnownCourses {
		courses = append(courses, v)
	}

	return courses
}

func GetKnownProfessors(faculty string) []string {
	var professors []string
	for _, v := range FacultyTimeTables[faculty].KnownProfessors {
		professors = append(professors, v)
	}

	return professors
}

func GetCurrentWeek() CurrentWeek {
	week := (((time.Now().UnixNano() / 1000000) - 1569794400000) / 604800000) + 1

	if week < 1 {
		week = 1
	}

	current := CurrentWeek{CurrentWeek: week}

	return current
}
