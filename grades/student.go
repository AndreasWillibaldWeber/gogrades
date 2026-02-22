package grades

import (
	"fmt"
	"strconv"

	"github.com/andreaswillibaldweber/gogrades/utilities"
)

type student struct {
	name    string
	matNr   string
	seatNr  string
	points  float64
	comment string
}

func NewStudent(name, matNr, seatNr string, points float64, comment string) *student {
	return &student{
		name:    name,
		matNr:   matNr,
		seatNr:  seatNr,
		points:  points,
		comment: comment,
	}
}

func (s student) Points() float64 {
	return s.points
}

func (s *student) Comment(comment string) {
	s.comment = comment
}

func (s student) String() string {
	return fmt.Sprintf(
		"Student{name: %q, matNr: %q, seatNr: %q, points: %.2f, comment: %q}",
		s.name, s.matNr, s.seatNr, s.points, s.comment,
	)
}

type students []student

func NewStudents() *students {
	return &students{}
}

func NewStudentsFromTable(table *utilities.Table) (*students, error) {
	var s students = *NewStudents()
	for _, row := range table.Rows() {
		name := fmt.Sprintf("%v", row[0])
		matNr := fmt.Sprintf("%v", row[1])
		seatNr := fmt.Sprintf("%v", row[2])
		points, err := strconv.ParseFloat(fmt.Sprintf("%v", row[3]), 64)
		if err != nil {
			return &students{}, fmt.Errorf("parse points for %q: %v", name, err)
		}
		comment := ""
		if len(row) > 4 {
			comment = fmt.Sprintf("%v", row[4])
		}
		s = *s.Add(NewStudent(name, matNr, seatNr, points, comment))
	}
	return &s, nil
}

func (s *students) Add(student *student) *students {
	*s = append(*s, *student)
	return s
}

func (s *students) Delete(matNr string) *students {
	for i, student := range *s {
		if student.matNr == matNr {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return s
		}
	}
	return s
}

func (s students) String() string {
	result := "Students:\n"
	for _, student := range s {
		result += fmt.Sprintf("- %s\n", student)
	}
	return result
}
