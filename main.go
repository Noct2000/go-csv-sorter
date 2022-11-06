package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	fmt.Println("== Started ==")
	inputFileName := flag.String("i", "", "Use a file with the name file-name as an input.")
	outputFileName := flag.String("o", "", "Use a file with the name file-name as an output.")
	sortingFilesIndex := flag.Int("f", 0, "Sort input lines by value number N.")
	isNotIgnoreHeader := flag.Bool("h", false, "The first line is a header that must be ignored during sorting but included in the output.")
	isReversedOrder := flag.Bool("r", false, "Sort input lines in reverse order.")
	flag.Parse()
	var content string
	if *inputFileName == "" {
		content = readFromConsole(*sortingFilesIndex, *isReversedOrder, *isNotIgnoreHeader)
	} else {
		content = readFromFile(*sortingFilesIndex, *isReversedOrder, *isNotIgnoreHeader, *inputFileName)
	}
	if content != "" {
		fmt.Println("sorted result:\n" + content)
		writeToFileIfPresent(content, *outputFileName)
	}
	fmt.Println("== Finished ==")

}

func readFromConsole(sortingFieldIndex int, isReversedOrder, isNotIgnoreHeader bool) string {
	scanner := bufio.NewScanner(os.Stdin)
	return processContent(sortingFieldIndex, isReversedOrder, isNotIgnoreHeader, scanner)
}

func readFromFile(sortingFieldIndex int, isReversedOrder, isNotIgnoreHeader bool, inputFile string) string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	content := processContent(sortingFieldIndex, isReversedOrder, isNotIgnoreHeader, fileScanner)
	return content
}

func processContent(sortingFieldIndex int, isReversedOrder, isNotIgnoreHeader bool, scanner *bufio.Scanner) string {
	var header string
	n := 0
	table := [][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ",")
		if n == 0 {
			n = len(row)
			if isNotIgnoreHeader {
				header = line
				continue
			}
		}
		if line == "" {
			break
		}
		if n != len(row) {
			log.Fatalf("Error: row has %d columns, but must have %d\n", len(row), n)
		}
		table = append(table, row)
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	sort.Slice(table, func(i, j int) bool {
		return compare(table[i][sortingFieldIndex], table[j][sortingFieldIndex], isReversedOrder)
	})
	var result strings.Builder
	if header != "" {
		result.WriteString(header)
		result.WriteString("\n")
	}
	for _, row := range table {
		result.WriteString(strings.Join(row, ","))
		result.WriteString("\n")
	}
	return result.String()
}

func writeToFileIfPresent(content, fileName string) {
	if fileName != "" {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func compare(first, next string, isReversed bool) bool {
	return first < next != isReversed
}
