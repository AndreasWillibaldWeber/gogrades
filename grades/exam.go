package grades

import (
	"fmt"
	"math"

	"github.com/andreaswillibaldweber/gogrades/utilities"
)

type exam struct {
	pMax     float64
	pPass    float64
	students students
}

func NewExam(pMax, pPass float64) exam {
	return exam{
		pMax:     pMax,
		pPass:    pPass,
		students: make(students, 0),
	}
}

func (e exam) Students() *students {
	return &e.students
}

func (e *exam) AddStudent(student *student) *students {
	e.students = append(e.students, *student)
	return &e.students
}

func (e *exam) AddStudents(students *students) *students {
	e.students = append(e.students, *students...)
	return &e.students
}

func (e *exam) DeleteStudent(matNr string) {
	for i, s := range e.students {
		if s.matNr == matNr {
			e.students = append(e.students[:i], e.students[i+1:]...)
			return
		}
	}
}

func (e exam) AmountStudents() int {
	return len(e.students)
}

func (e exam) Grade(s student) float64 {
	return e.LinearGrading(s.Points())
}

func (e exam) LinearGrading(points float64) float64 {

	if points < e.pPass {
		return 5.0
	}
	if points > e.pMax {
		return 1.0
	}

	return e.gRounded(points)
}

func (e exam) gRounded(p float64) float64 {
	return math.Floor(10*e.gStep(p)+0.5) / 10
}

func (e exam) graw(p float64) float64 {
	return 1 + 3*((e.pMax-p)/(e.pMax-e.pPass))
}

func (e exam) k(p float64) float64 {
	return math.Floor(3*(e.graw(p)-1) + 0.5)
}

func (e exam) gStep(p float64) float64 {
	return 1 + e.k(p)/3
}

func (e exam) GradingKeyTable() *utilities.Table {

	type grading struct {
		nr         int64
		points     float64
		percentage float64
		grade      float64
	}

	grades := make([]grading, 0)
	lastGrade := 0.0
	for i := 0.0; i < e.pMax+0.1; i += 0.5 {
		if lastGrade != e.LinearGrading(i) || i >= e.pMax-0.2 {
			grades = append(grades, grading{points: i, percentage: i / e.pMax * 100, grade: e.LinearGrading(i)})
			lastGrade = e.LinearGrading(i)
		}
	}

	header := []string{"Nr", "Points", "%", "Grade"}
	rows := make([]utilities.TableRow, 0)
	for nr, g := range grades {
		row := utilities.TableRow{nr, g.points, g.percentage, g.grade}
		rows = append(rows, row)
	}
	hooks := map[int]utilities.FormatHook{
		1: utilities.BuildDecimalFormatHook(1),
		2: utilities.BuildPercentageFormatHook(1),
		3: utilities.BuildDecimalFormatHook(1),
	}
	table := utilities.NewTable(header, rows)
	table.SetFormatHooks(hooks)
	return table
}

func (e exam) GradingKeyString() string {

	return fmt.Sprintf(
		"Grading key with %.2f points maximum and %.2f points passing:\n%s", e.pMax, e.pPass,
		e.GradingKeyTable().FormatTableRight([]int{1, 2}),
	)
}

func (e exam) GradedStudentTable() *utilities.Table {
	header := []string{"Student Name", "Mat", "Seat", "Points", "%", "Grade", "Comment"}
	rows := make([]utilities.TableRow, 0)
	for _, s := range e.students {
		row := utilities.TableRow{s.name, s.matNr, s.seatNr, s.points, 100 * s.points / e.pMax, e.Grade(s), s.comment}
		rows = append(rows, row)
	}
	hooks := map[int]utilities.FormatHook{
		3: utilities.BuildDecimalFormatHook(1),
		4: utilities.BuildPercentageFormatHook(1),
		5: utilities.BuildDecimalFormatHook(1),
	}
	table := utilities.NewTable(header, rows)
	table.SetFormatHooks(hooks)
	return table
}

func (e exam) GradedStudentString() string {
	return fmt.Sprintf(
		"Exam with %d students:\n%s", e.AmountStudents(),
		e.GradedStudentTable().FormatTableRight([]int{3, 4}),
	)
}

func (e exam) String() string {
	return e.GradingKeyString() + "\n" + e.GradedStudentString()
}
