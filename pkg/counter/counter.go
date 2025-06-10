package counter

import (
	"bufio"
	"os"
	"strings"
)

// CountLines opens a path and returns its line count.
func CountLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	n := 0
	for scanner.Scan() {
		n++
	}
	return n, scanner.Err()
}

// CountLinesAndWords counts lines and word frequencies in a file.
func CountLinesAndWords(path string) (int, map[string]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	lineCount := 0
	wordFreq := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		for _, word := range strings.Fields(line) {
			wordFreq[word]++
		}
	}
	return lineCount, wordFreq, scanner.Err()
}
