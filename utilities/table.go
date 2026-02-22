package utilities

import "fmt"

type FormatHook func(value float64) string

type Table struct {
	header      []string
	rows        []TableRow
	formatHooks map[int]FormatHook
}

type TableRow []any
type TableRows []TableRow

func NewEmptyTable(headers []string) *Table {
	return &Table{
		header: headers,
		rows:   make(TableRows, 0),
	}
}

func NewTable(headers []string, rows []TableRow) *Table {
	table := NewEmptyTable(headers).SetRows(rows)
	return table
}

func NewTableFromCSV(filepath string) (*Table, error) {
	table, err := ReadCSV(filepath)
	return table, err
}

func (t Table) Headers() []string {
	return t.header
}

func (t *Table) SetHeaders(header []string) Table {
	t.header = header
	return *t
}

func (t *Table) SetFormatHooks(hooks map[int]FormatHook) *Table {
	t.formatHooks = hooks
	return t
}

func (t *Table) AddFormatHook(colIdx int, hook FormatHook) *Table {
	t.formatHooks[colIdx] = hook
	return t
}

func (t *Table) ClearFormatHooks() {
	t.formatHooks = make(map[int]FormatHook)
}

func (t *Table) ClearHeaders() {
	t.header = []string{}
}

func (t Table) Rows() []TableRow {
	return t.rows
}

func (t *Table) SetRows(rows []TableRow) *Table {
	t.rows = rows
	return t
}

func (t *Table) AddRow(row TableRow) Table {
	t.rows = append(t.rows, row)
	return *t
}

func (t *Table) ClearRows() {
	t.rows = []TableRow{}
}

func (t *Table) DeleteRow(index int) {
	if index < 0 || index >= len(t.rows) {
		return
	}
	t.rows = append(t.rows[:index], t.rows[index+1:]...)
}

func (t Table) ToCSV(filepath string) error {
	return WriteCSV(filepath, t)
}

func (t Table) String() string {
	return t.FormatTable(nil)
}

func (t Table) FormatTable(rightAlignCols []int) string {
	if len(t.rows) == 0 {
		return fmt.Sprintf("Empty table (%d columns)\n", len(t.header))
	}

	widths := t.calculateWidths()
	separator := buildSeparator(widths)

	out := separator + "\n"
	out += t.buildHeaderRow(widths, rightAlignCols) + "\n"
	out += separator + "\n"
	out += t.buildDataRows(widths, rightAlignCols) + "\n"
	out += separator + "\n"

	return out
}

func (t Table) FormatTableRight(rightAlignCols []int) string {
	if len(t.rows) == 0 {
		return fmt.Sprintf("Empty table (%d columns)\n", len(t.header))
	}

	widths := t.calculateWidths()
	separator := buildSeparator(widths)

	out := separator + "\n"
	out += t.buildHeaderRow(widths, rightAlignCols) + "\n"
	out += separator + "\n"
	out += t.buildDataRows(widths, rightAlignCols) + "\n"
	out += separator + "\n"

	return out
}

func (t Table) calculateWidths() []int {
	widths := make([]int, len(t.header))

	for i, h := range t.header {
		widths[i] = len(h)
	}

	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) {
				widths[i] = maxInt(widths[i], len(fmt.Sprintf("%v", cell)))
			}
		}
	}

	return widths
}

func buildSeparator(widths []int) string {
	separator := "+"
	for _, w := range widths {
		separator += fmt.Sprintf("%s+", repeatString("-", w+2))
	}
	return separator
}

func (t Table) buildHeaderRow(widths []int, rightAlignCols []int) string {
	row := "|"
	for i, h := range t.header {
		row += " " + padCell(h, widths[i], isRightAligned(i, rightAlignCols)) + " |"
	}
	return row
}

func (t Table) buildDataRows(widths []int, rightAlignCols []int) string {
	out := ""
	for _, row := range t.rows {
		formatted := t.formatRowStrings(row)
		out += buildTableRowFromStrings(formatted, widths, rightAlignCols) + "\n"
	}
	return out[:len(out)-1]
}

func (t Table) formatRowStrings(row TableRow) []string {
	out := make([]string, len(row))
	for i, cell := range row {
		cellStr := fmt.Sprintf("%v", cell)
		if t.formatHooks != nil {
			if hook, ok := t.formatHooks[i]; ok && hook != nil {
				if f, okf := cell.(float64); okf {
					cellStr = hook(f)
				}
			}
		}
		out[i] = cellStr
	}
	return out
}

func buildTableRowFromStrings(formatted []string, widths []int, rightAlignCols []int) string {
	line := "|"
	for i, cellStr := range formatted {
		if i < len(widths) {
			line += " " + padCell(cellStr, widths[i], isRightAligned(i, rightAlignCols)) + " |"
		}
	}
	return line
}

func padCell(cell string, width int, rightAligned bool) string {
	if rightAligned {
		return fmt.Sprintf("%*s", width, cell)
	}
	return fmt.Sprintf("%-*s", width, cell)
}

func isRightAligned(colIdx int, rightAlignCols []int) bool {
	if rightAlignCols == nil {
		return false
	}
	for _, idx := range rightAlignCols {
		if idx == colIdx {
			return true
		}
	}
	return false
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func BuildDecimalFormatHook(decimal int) FormatHook {
	formatStr := fmt.Sprintf("%%.%df", decimal)
	return func(value float64) string {
		return fmt.Sprintf(formatStr, value)
	}
}

func BuildPercentageFormatHook(decimal int) FormatHook {
	formatStr := fmt.Sprintf("%%.%df%%%%", decimal)
	return func(value float64) string {
		return fmt.Sprintf(formatStr, value)
	}
}
