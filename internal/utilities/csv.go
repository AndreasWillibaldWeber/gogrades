package utilities

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadCSV(filepath string) (*Table, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return NewEmptyTable([]string{}), fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	return readCSVFromReader(f)
}

func readCSVFromReader(r io.Reader) (*Table, error) {
	reader := csv.NewReader(r)

	header, err := reader.Read()
	if err != nil {
		return NewEmptyTable([]string{}), fmt.Errorf("read header: %w", err)
	}

	table := NewEmptyTable(header)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return NewEmptyTable([]string{}), fmt.Errorf("read row: %w", err)
		}

		if len(row) < 4 {
			return NewEmptyTable([]string{}), fmt.Errorf("invalid row: expected at least 4 columns (name, matNr, seatNr, points), got %d", len(row))
		}

		points, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			return NewEmptyTable([]string{}), fmt.Errorf("parse points for %q: %w", row[0], err)
		}

		comment := ""
		if len(row) > 4 {
			comment = strings.TrimSpace(row[4])
		}

		table.AddRow(TableRow{strings.TrimSpace(row[0]), strings.TrimSpace(row[1]), strings.TrimSpace(row[2]), points, comment})
	}

	return table, nil
}

func WriteCSV(filepath string, table Table) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	return writeCSVToWriter(f, table)
}

func writeCSVToWriter(w io.Writer, table Table) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write(table.header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, row := range table.rows {
		stringRow := table.formatRowStrings(row)
		if err := writer.Write(stringRow); err != nil {
			return fmt.Errorf("write row for %q: %w", row[0], err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	return nil
}
