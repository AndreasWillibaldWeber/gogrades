package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andreaswillibaldweber/gogrades/internal/utilities"
)

type GUI struct {
	window fyne.Window

	pMax  float64
	pPass float64

	maxPointsEntry  *widget.Entry
	passPointsEntry *widget.Entry
	statusLabel     *widget.Label

	loadedCSVPath string
	loadedTable   *utilities.Table

	gradedStudents *utilities.Table
	gradingKey     *utilities.Table

	gradedTable *tableAdapter
	keyTable    *tableAdapter
}

func newGUI(pMax, pPass float64) *GUI {
	a := app.NewWithID("gogrades")
	w := a.NewWindow("GoGrades")
	w.Resize(fyne.NewSize(windowWidth, windowHeight))

	g := &GUI{
		window:          w,
		pMax:            pMax,
		pPass:           pPass,
		maxPointsEntry:  widget.NewEntry(),
		passPointsEntry: widget.NewEntry(),
		statusLabel:     widget.NewLabel("No CSV loaded. Use File -> Open CSV..."),
		gradedTable: newTableAdapter(func(colIdx int) bool {
			return colIdx == 3 || colIdx == 4 || colIdx == 5
		}),
		keyTable: newTableAdapter(func(colIdx int) bool {
			return colIdx == 1 || colIdx == 2 || colIdx == 3
		}),
	}

	g.maxPointsEntry.SetText(fmt.Sprintf("%.1f", pMax))
	g.passPointsEntry.SetText(fmt.Sprintf("%.1f", pPass))
	g.window.SetMainMenu(g.buildMainMenu())
	g.window.SetContent(g.buildContent())
	return g
}

func (g *GUI) buildControls() *fyne.Container {
	return container.NewHBox(
		widget.NewLabel("Max Points"),
		container.NewGridWrap(fyne.NewSize(90, g.maxPointsEntry.MinSize().Height), g.maxPointsEntry),
		widget.NewLabel("Pass Points"),
		container.NewGridWrap(fyne.NewSize(90, g.passPointsEntry.MinSize().Height), g.passPointsEntry),
		widget.NewButton("Apply", g.applySettings),
		g.statusLabel,
	)
}

func (g *GUI) buildContent() fyne.CanvasObject {
	leftPane := container.NewBorder(widget.NewLabel("Graded Students"), nil, nil, nil, g.gradedTable.Widget())
	rightPane := container.NewBorder(widget.NewLabel("Grading Key"), nil, nil, nil, g.keyTable.Widget())
	split := container.NewHSplit(leftPane, rightPane)
	split.Offset = splitOffset
	return container.NewBorder(g.buildControls(), nil, nil, nil, split)
}

func ShowExamTables(pMax, pPass float64, preloadCSVPath string) error {
	g := newGUI(pMax, pPass)

	if strings.TrimSpace(preloadCSVPath) != "" {
		if err := g.loadCSVPath(preloadCSVPath); err != nil {
			return err
		}
	}

	if err := g.rebuildTables(); err != nil {
		return err
	}
	g.renderTables()
	g.window.ShowAndRun()
	return nil
}
