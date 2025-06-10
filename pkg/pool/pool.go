package pool

import "sync"

type Result struct {
	Path  string
	Lines int
	Err   error
}

// Start spins up workers to apply countFn to each path.
// Emits Result on the returned channel.
func Start(paths <-chan string, workers int, countFn func(string) (int, error)) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range paths {
				n, err := countFn(p)
				out <- Result{Path: p, Lines: n, Err: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
