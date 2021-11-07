package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

const countTodo int = -1

type item struct {
	name  string
	count int
}

func main() {
	absolutePath, pathError := filepath.Abs(os.Args[1])
	if pathError != nil {
		log.Fatalf("Error reading file: %s\n", pathError)
	}

	categoryMapping, parseError := parseFile(absolutePath)
	if parseError != nil {
		log.Fatalf("Error parsing:\n%s\n", parseError)
	}

	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "Count\tName")
	fmt.Fprintln(writer, "----\t----")

	var overallCount int
	for category, items := range categoryMapping {
		var categoryCount int
		for _, item := range items {
			if item.count != countTodo {
				categoryCount += item.count
			}
		}
		fmt.Fprintf(writer, "%d\t%s\n", categoryCount, category)
		overallCount += categoryCount
	}

	writer.Flush()

	fmt.Printf("\n\nOverall Item count: %d\n", overallCount)
}

func parseFile(filePath string) (map[string][]item, error) {
	categoryMapping := make(map[string][]item)

	fileHandle, fileErr := os.Open(filePath)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	var lastCategory string
	lastLine := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "##") {
			// Category Heading
			lastCategory = strings.TrimSpace(strings.TrimPrefix(line, "##"))
		} else if line != "" {
			// Item
			cells := strings.Split(line, ",")
			if len(cells) != 2 {
				return nil, fmt.Errorf("column count in line %d is %d instead of %d: '%s'", lastLine, len(cells), 2, line)
			}

			cellTwo := strings.TrimSpace(cells[1])
			var count int
			if cellTwo == "TODO" {
				count = countTodo
			} else {
				int32Count, countErr := strconv.ParseInt(cellTwo, 10, 32)
				if countErr != nil {
					return nil, fmt.Errorf("malformed count on line %d: '%s'", lastLine, cellTwo)
				}
				count = int(int32Count)
			}

			//Skip for now
			if count == countTodo {
				continue
			}

			categoryMapping[lastCategory] = append(categoryMapping[lastCategory], item{
				name:  strings.TrimSpace(cells[0]),
				count: count,
			})
		}

		lastLine++
	}

	return categoryMapping, nil
}
