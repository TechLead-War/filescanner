package pool

import "sync"

func StartWithWords(paths <-chan string, workers int, countFn func(string) (int, map[string]int, error)) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range paths {
				lines, words, err := countFn(p)
				out <- Result{Path: p, Lines: lines, Words: words, Err: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
