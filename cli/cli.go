package cli

import (
	"flag"
	"fmt"
)

type flags struct {
	gstud   bool
	gkey    bool
	pmax    float64
	ppass   float64
	csvFile string
	saveCSV bool
}

func (f flags) GStud() bool {
	return f.gstud
}

func (f flags) GKey() bool {
	return f.gkey
}

func (f flags) PMax() float64 {
	return f.pmax
}

func (f flags) PPass() float64 {
	return f.ppass
}

func (f flags) CSVFile() string {
	return f.csvFile
}

func (f flags) SaveCSV() bool {
	return f.saveCSV
}

func (f flags) String() string {
	return fmt.Sprintf("pmax: %v, ppass: %v, csvFile: %s, saveCSV: %t, gkey: %t, gstud: %t", f.pmax, f.ppass, f.csvFile, f.saveCSV, f.gkey, f.gstud)
}

func ParseFlags() flags {
	gkey := flag.Bool("gkey", false, "show grading key")
	gstud := flag.Bool("gstud", false, "show student grading")
	pmax := flag.Float64("pmax", 90, "maximum points")
	ppass := flag.Float64("ppass", 45, "passing points")
	csvFile := flag.String("csvfile", "", "path to CSV file with student data")
	saveCSV := flag.Bool("savecsv", false, "path to save CSV file with student data (overwrites existing file)")

	flag.Parse()

	return flags{
		gstud:   *gstud,
		gkey:    *gkey,
		pmax:    *pmax,
		ppass:   *ppass,
		csvFile: *csvFile,
		saveCSV: *saveCSV,
	}
}
