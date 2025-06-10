package main

import (
	"filescanner/pkg/counter"
	"filescanner/pkg/pool"
	"filescanner/pkg/scanner"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

var cmds = map[string]*flag.FlagSet{}

func main() {
	// scan command and flags
	scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
	dir := scanCmd.String("dir", ".", "directory to scan")
	ext := scanCmd.String("ext", ".txt", "file extension")
	workers := scanCmd.Int("workers", 0, "number of goroutines (0=NumCPU)")
	mode := scanCmd.String("mode", "lines", "scan mode: 'lines' or 'words'")
	cmds["scan"] = scanCmd

	// command dispatch
	args := os.Args[1:]
	var cmdName string
	var cmdArgs []string

	switch {
	case len(args) == 0:
		cmdName, cmdArgs = "scan", []string{}
	case isHelp(args[0]):
		printUsage()
		os.Exit(0)
	case cmds[args[0]] != nil:
		cmdName, cmdArgs = args[0], args[1:]
	case strings.HasPrefix(args[0], "-"):
		cmdName, cmdArgs = "scan", args
	default:
		fmt.Printf("Unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}

	fs := cmds[cmdName]
	fs.Parse(cmdArgs)

	if cmdName == "scan" {
		runScan(*dir, *ext, *workers, *mode)
	}
}

func runScan(dir, ext string, workers int, mode string) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	paths := scanner.Walk(dir, ext)

	switch mode {
	case "lines":
		results := pool.Start(paths, workers, counter.CountLines)
		totalF, totalL := 0, 0
		for r := range results {
			if r.Err != nil {
				fmt.Fprintf(os.Stderr, "error %s: %v\n", r.Path, r.Err)
				continue
			}
			fmt.Printf("%s: %d\n", r.Path, r.Lines)
			totalF++
			totalL += r.Lines
		}
		fmt.Printf("\nScanned %d files, %d total lines\n", totalF, totalL)

	case "words":
		results := pool.StartWithWords(paths, workers, counter.CountLinesAndWords)
		totalF, totalL := 0, 0
		wordTotal := make(map[string]int)

		for r := range results {
			if r.Err != nil {
				fmt.Fprintf(os.Stderr, "error %s: %v\n", r.Path, r.Err)
				continue
			}
			fmt.Printf("%s: %d lines\n", r.Path, r.Lines)
			totalF++
			totalL += r.Lines
			for word, count := range r.Words {
				wordTotal[word] += count
			}
		}

		fmt.Printf("\nScanned %d files, %d total lines\n", totalF, totalL)

		type wordFreq struct {
			word  string
			count int
		}
		var list []wordFreq
		for word, count := range wordTotal {
			list = append(list, wordFreq{word, count})
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].count > list[j].count
		})
		fmt.Println("\nTop 10 words:")
		for i := 0; i < len(list) && i < 10; i++ {
			fmt.Printf("%s: %d\n", list[i].word, list[i].count)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s (use 'lines' or 'words')\n", mode)
		os.Exit(1)
	}
}

func isHelp(s string) bool {
	return s == "-h" || s == "--help" || s == "help"
}

func printUsage() {
	fmt.Printf("Usage: %s [command] [options]\n\n", os.Args[0])
	fmt.Println("Commands:")
	for name := range cmds {
		fmt.Printf("  %s\n", name)
	}
	fmt.Println("\nUse '<command> --help' to see its options.\n")
	for name, fs := range cmds {
		fmt.Printf("Options for '%s':\n", name)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Printf("  -%s\t%s (default %q)\n", f.Name, f.Usage, f.DefValue)
		})
		fmt.Println()
	}
}
