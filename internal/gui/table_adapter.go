package gui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/andreaswillibaldweber/gogrades/internal/utilities"
)

type cellFormatter func(colIdx int, value any) string

type tableAdapter struct {
	table      *widget.Table
	headers    []string
	rows       [][]string
	rightAlign func(colIdx int) bool
}

func newTableAdapter(rightAlign func(colIdx int) bool) *tableAdapter {
	t := &tableAdapter{headers: []string{}, rows: [][]string{}, rightAlign: rightAlign}
	t.table = t.newWidget()
	return t
}

func (t *tableAdapter) Widget() *widget.Table {
	return t.table
}

func (t *tableAdapter) newWidget() *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(t.rows), len(t.headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id.Row < 0 || id.Col < 0 || id.Row >= len(t.rows) || id.Col >= len(t.rows[id.Row]) {
				label.SetText("")
				return
			}
			label.SetText(t.rows[id.Row][id.Col])
			if t.rightAlign(id.Col) {
				label.Alignment = fyne.TextAlignTrailing
				return
			}
			label.Alignment = fyne.TextAlignLeading
		},
	)
	table.CreateHeader = t.createHeader
	table.UpdateHeader = t.updateHeader
	table.SetColumnWidth(-1, defaultRowHeader)
	return table
}

func (t *tableAdapter) createHeader() fyne.CanvasObject {
	label := widget.NewLabel("")
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

func (t *tableAdapter) updateHeader(id widget.TableCellID, obj fyne.CanvasObject) {
	label := obj.(*widget.Label)
	switch {
	case id.Row == -1 && id.Col >= 0:
		label.Alignment = fyne.TextAlignCenter
		if id.Col < len(t.headers) {
			label.SetText(t.headers[id.Col])
			return
		}
		label.SetText("")
	case id.Col == -1 && id.Row >= 0:
		label.Alignment = fyne.TextAlignTrailing
		label.SetText(strconv.Itoa(id.Row + 1))
	default:
		label.Alignment = fyne.TextAlignLeading
		label.SetText("")
	}
}

func (t *tableAdapter) setData(src *utilities.Table, formatter cellFormatter) {
	headers := src.Headers()
	rows := src.Rows()

	t.headers = make([]string, len(headers))
	copy(t.headers, headers)
	t.rows = make([][]string, len(rows))

	for rowIdx, row := range rows {
		formattedRow := make([]string, len(t.headers))
		for colIdx := range t.headers {
			if colIdx < len(row) {
				formattedRow[colIdx] = formatter(colIdx, row[colIdx])
			}
		}
		t.rows[rowIdx] = formattedRow
	}
}

func (t *tableAdapter) clear() {
	t.headers = []string{}
	t.rows = [][]string{}
	t.table.SetColumnWidth(-1, defaultRowHeader)
	t.table.Refresh()
}

func (t *tableAdapter) resizeToFit(availableWidth float32) {
	if len(t.headers) == 0 {
		t.table.SetColumnWidth(-1, defaultRowHeader)
		return
	}

	baseWidths, baseTotal := t.baseColumnWidths()
	rowHeaderWidth := t.rowHeaderWidth()
	dataAvailable := t.availableDataWidth(availableWidth, rowHeaderWidth)
	scale := t.scale(baseTotal, dataAvailable)

	for colIdx, baseWidth := range baseWidths {
		width := baseWidth * scale
		if width < 56 {
			width = 56
		}
		t.table.SetColumnWidth(colIdx, width)
	}
	t.table.SetColumnWidth(-1, rowHeaderWidth)
}

func (t *tableAdapter) baseColumnWidths() ([]float32, float32) {
	base := make([]float32, len(t.headers))
	total := float32(0)

	for colIdx, header := range t.headers {
		maxLen := len(header)
		for _, row := range t.rows {
			if colIdx < len(row) && len(row[colIdx]) > maxLen {
				maxLen = len(row[colIdx])
			}
		}
		width := float32(maxLen)*float32(approxCharWidth) + float32(cellPadding)
		if width < 90 {
			width = 90
		}
		if width > 320 {
			width = 320
		}
		base[colIdx] = width
		total += width
	}
	return base, total
}

func (t *tableAdapter) rowHeaderWidth() float32 {
	digits := len(strconv.Itoa(len(t.rows)))
	width := float32(digits)*float32(approxCharWidth) + 28
	if width < 48 {
		return 48
	}
	return width
}

func (t *tableAdapter) availableDataWidth(availableWidth, rowHeaderWidth float32) float32 {
	dataAvailable := availableWidth - rowHeaderWidth - 24
	minWidth := float32(len(t.headers)) * 56
	if dataAvailable < minWidth {
		return minWidth
	}
	return dataAvailable
}

func (t *tableAdapter) scale(baseTotal, dataAvailable float32) float32 {
	if baseTotal <= 0 || baseTotal <= dataAvailable {
		return 1.0
	}
	return dataAvailable / baseTotal
}

func formatGradedCell(colIdx int, value any) string {
	v, ok := value.(float64)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	if colIdx == 3 || colIdx == 5 {
		return fmt.Sprintf("%.1f", v)
	}
	if colIdx == 4 {
		return fmt.Sprintf("%.1f%%", v)
	}
	return fmt.Sprintf("%v", value)
}

func formatGradingKeyCell(colIdx int, value any) string {
	v, ok := value.(float64)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	if colIdx == 1 || colIdx == 3 {
		return fmt.Sprintf("%.1f", v)
	}
	if colIdx == 2 {
		return fmt.Sprintf("%.1f%%", v)
	}
	return fmt.Sprintf("%v", value)
}
