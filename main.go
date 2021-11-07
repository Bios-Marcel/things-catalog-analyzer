package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

const countTodo int = -1

type item struct {
	name  string
	count int
}

type category struct {
	name  string
	items []item
}

func (c *category) count() int {
	var count int
	for _, item := range c.items {
		count += item.count
	}
	return count
}

func main() {
	absolutePath, pathError := filepath.Abs(os.Args[1])
	if pathError != nil {
		log.Fatalf("Error reading file: %s\n", pathError)
	}

	categories, parseError := parseFile(absolutePath)
	if parseError != nil {
		log.Fatalf("Error parsing:\n%s\n", parseError)
	}

	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "Count\tName")
	fmt.Fprintln(writer, "----\t----")

	sort.Slice(categories, func(a, b int) bool {
		return categories[a].count() > categories[b].count()
	})

	var overallCount int
	for _, category := range categories {
		categoryCount := category.count()
		fmt.Fprintf(writer, "%d\t%s\n", categoryCount, category.name)
		overallCount += categoryCount
	}

	if flushError := writer.Flush(); flushError != nil {
		log.Fatalf("Error printing results: %s\n", flushError)
	}

	fmt.Printf("\n\nOverall Item count: %d\n", overallCount)
}

func parseFile(filePath string) ([]*category, error) {
	fileHandle, fileErr := os.Open(filePath)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	defer fileHandle.Close()

	var categories []*category
	var lastCategory *category
	lastLine := 1

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "##") {
			// Category Heading
			lastCategory = &category{
				name: strings.TrimSpace(strings.TrimPrefix(line, "##")),
			}
			categories = append(categories, lastCategory)
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

			lastCategory.items = append(lastCategory.items, item{
				name:  strings.TrimSpace(cells[0]),
				count: count,
			})
		}

		lastLine++
	}

	return categories, nil
}
