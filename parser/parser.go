package parser

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const SEPARATOR string = ";"

var firstDayOfWeek int64
var maxMinuteOfDay int64
var minMinuteOfDay int64

var schoolYear int64
var schoolYearStartDay int64
var schoolYearStartMonth int64

var segmentDuration int64

var semester1BeginDay int64
var semester1BeginMonth int64

var semester1EndDay int64
var semester1EndMonth int64

var serial string
var orgCity string
var orgName string

var startOfSchoolYear time.Time

var RoomsMap = make(map[int64]Room)
var GroupsMap = make(map[int64]Group)
var TutorsMap = make(map[int64]Tutor)
var CoursesMap = make(map[int64]Course)
var BranchesMap = make(map[int64]Branch)
var ProgramsMap = make(map[int64]Program)
var ExecTypesMap = make(map[int64]ExecType)
var TimeSlotsMap = make(map[int64]TimeSlot)

var Schedules = make([]Schedule, 0)
var reservations []Reservation

type Branch struct {
	Id          int64
	Name        string
	translation string
	ProgramId   int64
	Year        int64
}

type Course struct {
	Id          int64
	Name        string
	note        string
	translation string
}

type ExecType struct {
	ExecTypeId int64
	ExecType   string
}

type Group struct {
	Id            int64
	Name          string
	BranchId      int64
	note          string
	ParentGroupId int64
	password      string
}

type Program struct {
	id          int64
	Name        string
	translation string
	Years       int64
}

type Reservation struct {
	begin    int64
	comment  string
	duration int64
	fromDay  int64
	fromWeek int64
	groupIds [] int64
	roomIds  [] int64
	toDay    int64
	toWeek   int64
	tutorIds [] int64
}

type Room struct {
	Id       int64
	RoomName string
}

type Schedule struct {
	BeginWeek       int64
	EndWeek         int64
	DayOfWeek       int64
	BeginTimeSlotId int64
	Duration        int64
	RoomId          int64
	CourseId        int64
	ExecTypeId      int64
	GroupIds        [] int64
	TutorIds        [] int64
}

type TimeSlot struct {
	Id              int64
	StartTimeString string
	EndTimeString   string
	RealTimeHour    int64
	RealTimeMin     int64
}

type Tutor struct {
	Id        int64
	Firstname string
	Lastname  string
	note      string
	password  string
}

func Parse() {
	file, err := os.Open("timetable.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := bufio.NewScanner(file)

	serial = readLine(data)
	orgCity = readLine(data)
	orgName = readLine(data)
	segmentDuration, _ = strconv.ParseInt(readLine(data), 10, 64)

	temp, _ := strconv.ParseInt(readLine(data), 10, 64)

	switch temp {
	case 0:
		firstDayOfWeek = 0
	case 1:
		firstDayOfWeek = 6
	case 2:
		firstDayOfWeek = 5
	default:
		firstDayOfWeek = 0
	}

	sSplit := strings.Split(readLine(data), SEPARATOR)

	schoolYear, _ = strconv.ParseInt(sSplit[0], 10, 64)
	schoolYearStartDay, _ = strconv.ParseInt(sSplit[1], 10, 64)
	schoolYearStartMonth, _ = strconv.ParseInt(sSplit[2], 10, 64)

	sSplit = strings.Split(readLine(data), SEPARATOR)

	semester1BeginDay, _ = strconv.ParseInt(sSplit[0], 10, 64)
	semester1BeginMonth, _ = strconv.ParseInt(sSplit[1], 10, 64)
	semester1EndDay, _ = strconv.ParseInt(sSplit[2], 10, 64)
	semester1EndMonth, _ = strconv.ParseInt(sSplit[3], 10, 64)

	var i int64

	i = 0
	n, _ := strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		slot := makeTimeSlot(i, readLine(data))
		TimeSlotsMap[slot.Id] = slot
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		room := makeRoom(readLine(data))
		RoomsMap[room.Id] = room
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		course := makeCourse(readLine(data))
		CoursesMap[course.Id] = course
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		execType := makeExecType(readLine(data))
		ExecTypesMap[execType.ExecTypeId] = execType
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		program := makeProgram(readLine(data))
		ProgramsMap[program.id] = program
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		branch := makeBranch(readLine(data))
		BranchesMap[branch.Id] = branch
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		tutor := makeTutor(readLine(data))
		TutorsMap[tutor.Id] = tutor
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		group := makeGroup(readLine(data))
		GroupsMap[group.Id] = group
	}

	i = 0
	n, _ = strconv.ParseInt(readLine(data), 10, 64)
	for ; i < n; i++ {
		line := readLine(data)
		Schedules = append(Schedules, makeSchedule(line))
	}

	if err := data.Err(); err != nil {
		log.Fatal(err)
	}
}

func readLine(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return scanner.Text()
	} else {
		return ""
	}
}

func makeTimeSlot(id int64, s string) TimeSlot {
	var realTimeHour int64
	var realTimeMin int64

	lSplit := strings.Split(s, SEPARATOR)
	startTimeString := lSplit[0]
	endTimeString := lSplit[1]

	possibleHour, _ := strconv.ParseInt(lSplit[2], 10, 64)
	possibleMinute, _ := strconv.ParseInt(lSplit[3], 10, 64)

	if possibleHour < 0 || possibleHour > 60 || possibleMinute < 0 || possibleMinute > 60 {
		startTimeSplit := strings.Split(startTimeString, ":")
		realTimeHour, _ = strconv.ParseInt(startTimeSplit[0], 10, 64)
		realTimeMin, _ = strconv.ParseInt(startTimeSplit[1], 10, 64)
	} else {
		realTimeHour = possibleHour
		realTimeMin = possibleMinute
	}
	return TimeSlot{id, startTimeString, endTimeString, realTimeHour, realTimeMin}
}

func makeRoom(s string) Room {
	sSplit := strings.Split(s, SEPARATOR)
	id, _ := strconv.ParseInt(sSplit[0], 10, 64)
	return Room{id, sSplit[1]}
}

func makeCourse(s string) Course {
	sSplit := strings.Split(s, SEPARATOR)
	note := ""
	translation := ""
	if len(sSplit) != 2 {
		if len(sSplit) == 3 {
			note = sSplit[2]
		} else if len(sSplit) == 4 {
			if sSplit[2] != "" {
				note = sSplit[2]
			}
			translation = sSplit[3]
		}
	}
	id, _ := strconv.ParseInt(sSplit[0], 10, 64)
	return Course{id, sSplit[1], note, translation}
}

func makeExecType(s string) ExecType {
	sSplit := strings.Split(s, SEPARATOR)

	id, _ := strconv.ParseInt(sSplit[0], 10, 64)
	return ExecType{id, sSplit[1]}
}

func makeProgram(s string) Program {
	sSplit := strings.Split(s, SEPARATOR)
	id, _ := strconv.ParseInt(sSplit[0], 10, 64)
	years, _ := strconv.ParseInt(sSplit[3], 10, 64)

	programName := sSplit[1]
	translation := ""

	if sSplit[2] != "" {
		translation = sSplit[2]
	}

	return Program{id, programName, translation, years}
}

func makeBranch(s string) Branch {
	sSplit := strings.Split(s, SEPARATOR)
	id, _ := strconv.ParseInt(sSplit[0], 10, 64)
	year, _ := strconv.ParseInt(sSplit[3], 10, 64)
	programId, _ := strconv.ParseInt(sSplit[4], 10, 64)

	translation := ""
	if sSplit[2] != "" {
		translation = sSplit[2]
	}

	return Branch{id, sSplit[1], translation, programId, year}
}

func makeTutor(s string) Tutor {
	sSplit := strings.Split(s, SEPARATOR)
	id, _ := strconv.ParseInt(sSplit[0], 10, 64)

	note := ""
	password := ""
	if len(sSplit) != 3 {
		if len(sSplit) == 4 {
			note = sSplit[3]
		} else if len(sSplit) == 5 {
			if sSplit[3] != "" {
				note = sSplit[3]
			}
			password = sSplit[4]
		}
	}

	return Tutor{id, sSplit[1], sSplit[2], note, password}
}

func makeGroup(s string) Group {
	sSplit := strings.Split(s, SEPARATOR)
	id, _ := strconv.ParseInt(sSplit[0], 10, 64)

	note := ""
	password := ""
	var parentGroupId int64 = 0
	branchId, _ := strconv.ParseInt(sSplit[2], 10, 64)

	if len(sSplit) == 4 {
		note = sSplit[3]
	} else if len(sSplit) == 5 {
		if sSplit[3] != "" {
			note = sSplit[3]
		}
		password = sSplit[4]
	} else if len(sSplit) == 6 {
		if sSplit[3] != "" {
			note = sSplit[3]
		}
		if sSplit[4] != "" {
			password = sSplit[4]
		}

		var err error

		parentGroupId, err = strconv.ParseInt(sSplit[5], 10, 64)

		if err != nil {
			parentGroupId = -1
		}
	}

	return Group{id, sSplit[1], branchId, note, parentGroupId, password}
}

func makeSchedule(s string) Schedule {
	var i int64
	var endMinuteOfDay int64

	sSplit := strings.Split(s, SEPARATOR)

	var groupIds = make([]int64, 0)

	numOfGroups, _ := strconv.ParseInt(sSplit[8], 10, 64)

	i = 9
	for ; i < numOfGroups+9; i++ {

		index, _ := strconv.ParseInt(sSplit[i], 10, 64)

		if GroupsMap[index].Name != "" && GroupsMap[index].Id != 0 {
			groupIds = append(groupIds, index)
		}
	}

	var tutorIds = make([]int64, 0)

	i = (numOfGroups + 9) + 1

	for ; i < int64(len(sSplit)); i++ {
		tutorId, _ := strconv.ParseInt(sSplit[i], 10, 64)
		if TutorsMap[tutorId].Firstname != "" && TutorsMap[tutorId].Id != 0 {
			tutorIds = append(tutorIds, tutorId)
		}
	}

	beginTimeSlotId, _ := strconv.ParseInt(sSplit[3], 10, 64)

	endTimeSlotId, _ := strconv.ParseInt(sSplit[4], 10, 64)

	endTimeSlotId = beginTimeSlotId + endTimeSlotId

	beginMinuteOfDay := (TimeSlotsMap[beginTimeSlotId].RealTimeHour * 60) + (TimeSlotsMap[beginTimeSlotId].RealTimeMin)

	if endTimeSlotId < int64(len(TimeSlotsMap)) {
		endMinuteOfDay = (TimeSlotsMap[endTimeSlotId].RealTimeHour * 60) + (TimeSlotsMap[endTimeSlotId].RealTimeMin)
	} else {
		endMinuteOfDay = maxMinuteOfDay - (segmentDuration / 2)
	}

	beginWeek, _ := strconv.ParseInt(sSplit[0], 10, 64)
	endWeek, _ := strconv.ParseInt(sSplit[1], 10, 64)
	dayOfWeek, _ := strconv.ParseInt(sSplit[2], 10, 64)
	btsid, _ := strconv.ParseInt(sSplit[3], 10, 64)
	duration := endMinuteOfDay - beginMinuteOfDay
	roomId, _ := strconv.ParseInt(sSplit[5], 10, 64)
	courseId, _ := strconv.ParseInt(sSplit[6], 10, 64)
	execId, _ := strconv.ParseInt(sSplit[7], 10, 64)

	return Schedule{beginWeek, endWeek, dayOfWeek, btsid,
		duration, roomId, courseId, execId, groupIds, tutorIds}
}

func Clean() {
	RoomsMap = make(map[int64]Room)
	GroupsMap = make(map[int64]Group)
	TutorsMap = make(map[int64]Tutor)
	CoursesMap = make(map[int64]Course)
	BranchesMap = make(map[int64]Branch)
	ProgramsMap = make(map[int64]Program)
	ExecTypesMap = make(map[int64]ExecType)
	TimeSlotsMap = make(map[int64]TimeSlot)
	Schedules = make([]Schedule, 0)
}
