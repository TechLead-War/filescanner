package counter

import (
	"bufio"
	"os"
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
