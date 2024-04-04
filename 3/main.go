package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type customSort struct {
	lines         []string
	key           int
	numeric       bool
	reverse       bool
	unique        bool
	byMonth       bool
	ignoreBlanks  bool
	checkSort     bool
	numericSuffix bool
}

func (cs *customSort) Len() int {
	return len(cs.lines)
}

func (cs *customSort) Swap(i, j int) {
	cs.lines[i], cs.lines[j] = cs.lines[j], cs.lines[i]
}

func (cs *customSort) Less(i, j int) bool {
	line1 := getColumnValue(cs.lines[i], cs.key)
	line2 := getColumnValue(cs.lines[j], cs.key)
	if cs.ignoreBlanks {
		line1 = strings.TrimSpace(line1)
		line2 = strings.TrimSpace(line2)
	}

	if cs.byMonth {
		line1 = convertToMonthName(line1)
		line2 = convertToMonthName(line2)
	}

	if cs.numericSuffix {
		line1 = convertToNumeric(convertToNumericSuffix(line1))
		line2 = convertToNumeric(convertToNumericSuffix(line2))
	}

	if cs.numeric {
		num1, err1 := strconv.Atoi(line1)
		num2, err2 := strconv.Atoi(line2)

		if err1 == nil && err2 == nil {
			line1 = fmt.Sprintf("%064d", num1)
			line2 = fmt.Sprintf("%064d", num2)
		}
	}

	if cs.reverse {

		if cs.checkSort && line1 < line2 {
			log.Fatal("disorder: ", line1)
		}

		return line1 > line2
	}

	if cs.checkSort && line1 < line2 {
		log.Fatal("disorder: ", line1)
	}

	return line1 < line2
}

func getColumnValue(line string, key int) string {
	cols := strings.Fields(line)
	if key > 0 && key <= len(cols) {
		return cols[key-1]
	}

	return line
}

func convertToMonthName(line string) string {
	fiels := strings.Fields(line)

	for _, field := range fiels {
		t, err := time.Parse("Jan", field)

		if err == nil {
			return fmt.Sprintf("%02d", t.Month())
		}
	}

	return line

}

func convertToNumericSuffix(line string) string {
	if strings.HasSuffix(line, "K") {
		num, _ := strconv.Atoi(strings.TrimSuffix(line, "K"))
		return fmt.Sprintf("%064d", num*1_000)
	}

	if strings.HasSuffix(line, "M") {
		num, _ := strconv.Atoi(strings.TrimSuffix(line, "M"))
		return fmt.Sprintf("%064d", num*1_000_000)
	}

	if strings.HasSuffix(line, "G") {
		num, _ := strconv.Atoi(strings.TrimSuffix(line, "G"))
		return fmt.Sprintf("%064d", num*1_000_000_000)
	}

	return line
}

func convertToNumeric(line string) string {
	num, err := strconv.Atoi(line)

	if err != nil {
		return line
	}
	return fmt.Sprintf("%064d", num)
}

func sortLines(cs *customSort) []string {
	if cs.unique {
		cs.lines = newSet(cs.lines)
	}

	sort.Sort(cs)
	return cs.lines
}

func newSet(lines []string) []string {
	m := make(map[string]struct{})
	var res []string

	for _, line := range lines {
		if _, ok := m[line]; !ok {
			m[line] = struct{}{}
			res = append(res, line)
		}
	}

	return res
}

func main() {
	key := flag.Int("k", -1, "Column sort by (1-indexed)")
	numeric := flag.Bool("n", false, "Sort numerically")
	reverse := flag.Bool("r", false, "Sort in reversed order")
	byMonth := flag.Bool("M", false, "Sort by month")
	ignoreBlanks := flag.Bool("b", false, "Ignore trailing spaces")
	checkSort := flag.Bool("c", false, "Check if the input is already sorted")
	numericSuffix := flag.Bool("h", false, "Compare human-readable numbers (e.g. 2K, 10M)")
	unique := flag.Bool("u", false, "Suppress lines that appear more than once")

	flag.Parse()
	inputFile := os.Args[len(os.Args)-1]

	inputFileHandle, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer inputFileHandle.Close()

	var lines []string
	scanner := bufio.NewScanner(inputFileHandle)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading input file:", err)
	}

	cs := &customSort{
		lines:         lines,
		key:           *key,
		numeric:       *numeric,
		reverse:       *reverse,
		byMonth:       *byMonth,
		ignoreBlanks:  *ignoreBlanks,
		checkSort:     *checkSort,
		numericSuffix: *numericSuffix,
		unique:        *unique,
	}

	sortedLines := sortLines(cs)

	outputFileHandle, err := os.Create("out.txt")
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outputFileHandle.Close()

	writer := bufio.NewWriter(outputFileHandle)

	for _, line := range sortedLines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal("Error writing to output file:", err)
		}
	}

	writer.Flush()

	log.Println("Sort completed")
}
