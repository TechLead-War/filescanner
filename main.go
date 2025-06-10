package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	root := "."
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && filepath.Ext(path) == ".txt" {
			count := countLines(path)
			fmt.Printf("%s: %d\n", path, count)
		}
		return nil
	})
}

func countLines(path string) int {
	f, _ := os.Open(path)
	defer f.Close()
	s := bufio.NewScanner(f)
	n := 0
	for s.Scan() {
		n++
	}
	return n
}
