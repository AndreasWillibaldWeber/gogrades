package gui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/andreaswillibaldweber/gogrades/internal/grades"
	"github.com/andreaswillibaldweber/gogrades/internal/utilities"
)

func (g *GUI) renderTables() {
	if g.gradedStudents == nil || g.gradingKey == nil {
		g.gradedTable.clear()
		g.keyTable.clear()
		return
	}

	g.gradedTable.setData(g.gradedStudents, formatGradedCell)
	g.keyTable.setData(g.gradingKey, formatGradingKeyCell)

	leftWidth, rightWidth := g.tablePaneWidths()
	g.gradedTable.resizeToFit(leftWidth)
	g.keyTable.resizeToFit(rightWidth)
	g.gradedTable.Widget().Refresh()
	g.keyTable.Widget().Refresh()
}

func (g *GUI) tablePaneWidths() (float32, float32) {
	canvasWidth := g.window.Canvas().Size().Width
	if canvasWidth <= 0 {
		canvasWidth = windowWidth
	}
	left := canvasWidth*float32(splitOffset) - 24
	right := canvasWidth*(1-float32(splitOffset)) - 24
	if left < 240 {
		left = 240
	}
	if right < 240 {
		right = 240
	}
	return left, right
}

func (g *GUI) rebuildTables() error {
	exam := grades.NewExam(g.pMax, g.pPass)
	if g.loadedTable != nil {
		students, err := grades.NewStudentsFromTable(g.loadedTable)
		if err != nil {
			return err
		}
		exam.AddStudents(students)
	}
	g.gradedStudents = exam.GradedStudentTable()
	g.gradingKey = exam.GradingKeyTable()
	return nil
}

func (g *GUI) applySettings() {
	maxVal, passVal, err := g.parseSettings()
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	g.pMax = maxVal
	g.pPass = passVal
	if err := g.rebuildTables(); err != nil {
		dialog.ShowError(fmt.Errorf("rebuild tables: %w", err), g.window)
		return
	}
	g.renderTables()
	g.statusLabel.SetText(fmt.Sprintf("Settings applied (max: %.1f, pass: %.1f)", g.pMax, g.pPass))
}

func (g *GUI) parseSettings() (float64, float64, error) {
	maxVal, err := strconv.ParseFloat(strings.TrimSpace(g.maxPointsEntry.Text), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid max points: %w", err)
	}
	passVal, err := strconv.ParseFloat(strings.TrimSpace(g.passPointsEntry.Text), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid pass points: %w", err)
	}
	if maxVal <= 0 {
		return 0, 0, fmt.Errorf("max points must be > 0")
	}
	if passVal < 0 {
		return 0, 0, fmt.Errorf("pass points must be >= 0")
	}
	if passVal >= maxVal {
		return 0, 0, fmt.Errorf("pass points must be smaller than max points")
	}
	return maxVal, passVal, nil
}

func (g *GUI) openCSVDialog() {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(fmt.Errorf("open file dialog: %w", err), g.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		uri := reader.URI()
		if uri == nil {
			dialog.ShowError(fmt.Errorf("could not resolve selected file"), g.window)
			return
		}
		if err := g.loadCSVPath(uri.Path()); err != nil {
			dialog.ShowError(err, g.window)
			return
		}
	}, g.window)
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	fileDialog.Show()
}

func (g *GUI) loadCSVPath(path string) error {
	table, err := utilities.NewTableFromCSV(path)
	if err != nil {
		return fmt.Errorf("read CSV: %w", err)
	}
	g.loadedTable = table
	g.loadedCSVPath = path
	if err := g.rebuildTables(); err != nil {
		return fmt.Errorf("parse students: %w", err)
	}
	g.renderTables()
	g.statusLabel.SetText(fmt.Sprintf("Loaded %s", g.loadedCSVPath))
	return nil
}

func (g *GUI) saveCSV() {
	if g.gradedStudents == nil || g.gradingKey == nil {
		dialog.ShowError(fmt.Errorf("no data loaded to save"), g.window)
		return
	}
	if strings.TrimSpace(g.loadedCSVPath) == "" {
		dialog.ShowError(fmt.Errorf("missing source CSV path for predefined save names"), g.window)
		return
	}

	basePath := strings.TrimSuffix(g.loadedCSVPath, filepath.Ext(g.loadedCSVPath))
	gradingKeyPath := basePath + "-grading-key.csv"
	gradedStudentsPath := basePath + "-graded.csv"

	errGradingKey := g.gradingKey.ToCSV(gradingKeyPath)
	errGradedStudents := g.gradedStudents.ToCSV(gradedStudentsPath)
	if errGradingKey != nil || errGradedStudents != nil {
		dialog.ShowError(fmt.Errorf("save CSV files (grading key: %v, graded students: %v)", errGradingKey, errGradedStudents), g.window)
		return
	}
	g.statusLabel.SetText(fmt.Sprintf("Saved %s and %s", gradingKeyPath, gradedStudentsPath))
}
