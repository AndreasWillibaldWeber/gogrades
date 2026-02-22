package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/andreaswillibaldweber/gogrades/cli"
	"github.com/andreaswillibaldweber/gogrades/grades"
	"github.com/andreaswillibaldweber/gogrades/utilities"
)

func Run() {
	flags := cli.ParseFlags()
	fmt.Printf("Flags>> %s \n\n", flags)

	exam := grades.NewExam(flags.PMax(), flags.PPass())

	if flags.GKey() {
		fmt.Println(exam.GradingKeyString())
	}

	table, err := utilities.NewTableFromCSV(flags.CSVFile())
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}
	students, err := grades.NewStudentsFromTable(table)
	if err != nil {
		fmt.Printf("Error parsing students from table: %v\n", err)
		return
	}
	exam.AddStudents(students)

	if flags.GStud() {
		fmt.Println(exam.GradedStudentString())
	}

	if flags.SaveCSV() {
		newpathGradingKey := strings.TrimSuffix(flags.CSVFile(), filepath.Ext(flags.CSVFile())) + "-grading-key.csv"
		newpathGradedStudent := strings.TrimSuffix(flags.CSVFile(), filepath.Ext(flags.CSVFile())) + "-graded.csv"

		err1 := exam.GradingKeyTable().ToCSV(newpathGradingKey)
		err2 := exam.GradedStudentTable().ToCSV(newpathGradedStudent)
		if err1 != nil || err2 != nil {
			fmt.Printf("Error writing CSV: %v\n", err)
			return
		}
		fmt.Printf("Exam data saved as CSV files.\n")
	}
}
