package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fields := flag.String("f", "-1", "Select fields (columns)")
	delimiter := flag.String("d", "\t", "Use different delimiter")
	separated := flag.Bool("s", false, "Only separated strings")

	flag.Parse()

	f, err := getFields(*fields)
	if err != nil {
		log.Fatal(err)
	}

	cut(f, *delimiter, *separated)
}

/*
	пример:
	```
		echo -e "1, 2, 3\n4, 5, 6\n" | go run main.go -f 1,3 -d ","
	```
*/

func getFields(f string) ([]int, error) {
	var res []int
	split := strings.Split(f, ",")

	for _, str := range split {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}

		res = append(res, num)
	}

	return res, nil
}

func cut(fields []int, delimiter string, separated bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, delimiter)
		if separated && !strings.Contains(line, delimiter) {
			continue
		}
		var selectedFields []string
		for _, field := range fields {
			if field < 0 {
				selectedFields = append(selectedFields, line)
				continue
			}
			if field <= len(parts) {
				selectedFields = append(selectedFields, parts[field-1])
			}
		}
		fmt.Println(strings.Join(selectedFields, delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading standard input:", err)
	}
}
