package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type FacultyWISE struct {
	ShortName string
	LongName  string
	Hash      string
}

type Faculty struct {
	ShortName string
	LongName  string
}

var FacultyMap = make(map[string]FacultyWISE)

func ParseFaculty() {
	file, err := os.Open("orgs.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := bufio.NewScanner(file)

	for ; data.Scan(); {
		text := data.Text()

		sSplit := strings.Split(text, SEPARATOR)

		fac := FacultyWISE{sSplit[0], sSplit[1], sSplit[2]}

		FacultyMap[fac.ShortName] = fac
	}
}

func GetKnownFaculties() []Faculty {
	var faculties []Faculty
	for _, v := range FacultyMap {
		faculty := Faculty{v.ShortName, v.LongName}
		faculties = append(faculties, faculty)
	}

	return faculties
}
