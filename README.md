# GoGrades

Simple Go project to compute grades from achieved points.
It load student information and points from a csv file and writes the result back to new csv files.

Build:

```bash
go build -o gogrades cmd/main.go 
go run cmd/main.go
```

Run examples:

```bash
./gogrades --csvfile example/students.csv --pmax 90 --ppass 45 --gkey --gstud --savecsv
```

Flags:
- `--gkey` show grading key
- `--gstud` show student grading
- `--pmax` maximum points (default 90)
- `--ppass` passing points (default 45)
- `--csvfile` path to CSV file with student data
- `--savecsv` save CSV file with student data to `csvfilepath-graded.csv` and grading key to `csvfilepath-grading-key.csv` (overwrites existing files)
