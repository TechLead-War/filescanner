package scanner

import (
	"io/fs"
	"path/filepath"
)

// Walk returns a channel of file paths under dir with the given ext.
func Walk(dir, ext string) <-chan string {
	out := make(chan string)
	go func() {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && filepath.Ext(path) == ext {
				out <- path
			}
			return nil
		})
		close(out)
	}()
	return out
}
