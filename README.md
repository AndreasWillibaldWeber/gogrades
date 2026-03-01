# GoGrades

Simple Go project to compute grades from achieved points.
It loads student information and points from a CSV file, then writes the results back to new CSV files.

Build the command-line tool:

```bash
mkdir ./build
go build -o ./gogrades/main.go
```

Run the command-line tool:

```bash
go run /gogrades/main.go
```

Installing the command-line tool:

1. Check your Go installation: [https://go.dev/doc/tutorial/compile-install](https://go.dev/doc/tutorial/compile-install)
    1. Add binary path to $PATH e.g. add to the file ~/.profile the command export PATH=$PATH:~/go/bin
    2. Set GOBIN path, e.g. with the command go env -w GOBIN=~/go/bin
   
2. Check repository structure: [https://go.dev/doc/modules/layout](https://go.dev/doc/modules/layout)
    1. Check the repository for version and release tags. Can be requested by @latest @v1.0.3 suffixes.

```bash
go install github.com/andreaswillibaldweber/gogrades/gogrades
```

Run examples:

```bash
./gogrades --csvfile example/students.csv --pmax 90 --ppass 45 --gkey --gstud --savecsv
```

Flags:
- `--gkey` show grading key
- `--gstud` show graded students
- `--gui` show GUI with graded students and grading key view
- `--pmax` maximum points (default 90)
- `--ppass` passing points (default 45)
- `--csvfile` path to CSV file with student data
- `--savecsv` save CSV file with graded students to `csvfilepath-graded.csv` and grading key to `csvfilepath-grading-key.csv` (overwrites existing files)

# Input format

Student table:

| Name          | Mat-Nr | Seat-Nr | Points | Comment           |
| ------------- | -----: | :-----: | -----: | ----------------- |
| Alice Johnson |  12001 |    A1   |   87.5 | Good performance  |
| Bob Smith     |  12002 |    A2   |     45 | Passing grade     |
| Charlie Brown |  12003 |    A3   |     92 | Excellent work    |
| Diana Prince  |  12004 |    B1   |   38.5 | Below passing     |
| Eve Davis     |  12005 |    B2   |     90 | Outstanding       |
| Frank Miller  |  12006 |    B3   |   55.5 | Satisfactory      |
| Grace Lee     |  12007 |    C1   |     78 | Solid performance |
| Henry Chen    |  12008 |    C2   |     42 | Just below pass   |
| Iris Wong     |  12009 |    C3   |   88.5 | Very good         |
| Jack Wilson   |  12010 |    D1   |     50 | Acceptable        |

Example student table (students.csv):
```csv
Name,Mat-Nr,Seat-Nr,Points,Comment
Alice Johnson,12001,A1,87.5,Good performance
Bob Smith,12002,A2,45,Passing grade
Charlie Brown,12003,A3,92,Excellent work
Diana Prince,12004,B1,38.5,Below passing
Eve Davis,12005,B2,90,Outstanding
Frank Miller,12006,B3,55.5,Satisfactory
Grace Lee,12007,C1,78,Solid performance
Henry Chen,12008,C2,42,Just below pass
Iris Wong,12009,C3,88.5,Very good
Jack Wilson,12010,D1,50,Acceptable
```

# Output format

Grading key table:

| Nr | Points |    %   | Grade |
| -: | -----: | :----: | ----: |
|  0 |    0.0 |  0.0%  |   5.0 |
|  1 |   45.0 |  50.0% |   4.0 |
|  2 |   48.0 |  53.3% |   3.7 |
|  3 |   53.0 |  58.9% |   3.3 |
|  4 |   58.0 |  64.4% |   3.0 |
|  5 |   63.0 |  70.0% |   2.7 |
|  6 |   68.0 |  75.6% |   2.3 |
|  7 |   73.0 |  81.1% |   2.0 |
|  8 |   78.0 |  86.7% |   1.7 |
|  9 |   83.0 |  92.2% |   1.3 |
| 10 |   88.0 |  97.8% |   1.0 |
| 11 |   90.0 | 100.0% |   1.0 |

Graded student table:

| Student Name  |   Mat | Seat | Points |    %   | Grade | Comment           |
| ------------- | ----: | :--: | -----: | :----: | ----: | ----------------- |
| Alice Johnson | 12001 |  A1  |   87.5 |  97.2% |   1.3 | Good performance  |
| Bob Smith     | 12002 |  A2  |   45.0 |  50.0% |   4.0 | Passing grade     |
| Charlie Brown | 12003 |  A3  |   92.0 | 102.2% |   1.0 | Excellent work    |
| Diana Prince  | 12004 |  B1  |   38.5 |  42.8% |   5.0 | Below passing     |
| Eve Davis     | 12005 |  B2  |   90.0 | 100.0% |   1.0 | Outstanding       |
| Frank Miller  | 12006 |  B3  |   55.5 |  61.7% |   3.3 | Satisfactory      |
| Grace Lee     | 12007 |  C1  |   78.0 |  86.7% |   1.7 | Solid performance |
| Henry Chen    | 12008 |  C2  |   42.0 |  46.7% |   5.0 | Just below pass   |
| Iris Wong     | 12009 |  C3  |   88.5 |  98.3% |   1.0 | Very good         |
| Jack Wilson   | 12010 |  D1  |   50.0 |  55.6% |   3.7 | Acceptable        |
