package config

import (
	"flag"
	"runtime"
)

type Config struct {
	Dir     string
	Ext     string
	Workers int
}

func Parse() Config {
	var c Config
	flag.StringVar(&c.Dir, "dir", ".", "directory to scan")
	flag.StringVar(&c.Ext, "ext", ".txt", "file extension to include")
	flag.IntVar(&c.Workers, "workers", 0, "number of goroutines (0=CPU count)")
	flag.Parse()
	if c.Workers <= 0 {
		c.Workers = runtime.NumCPU()
	}
	return c
}
