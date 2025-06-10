package main

import (
	"filescanner/pkg/counter"
	"filescanner/pkg/pool"
	"filescanner/pkg/scanner"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var cmds = map[string]*flag.FlagSet{} // Global so printUsage() can access

func main() {
	// Define "scan" subcommand and flags
	scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
	dir := scanCmd.String("dir", ".", "directory to scan")
	ext := scanCmd.String("ext", ".txt", "file extension")
	workers := scanCmd.Int("workers", 0, "number of goroutines (0=NumCPU)")
	cmds["scan"] = scanCmd

	// Determine command
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

	// Parse command-specific flags
	fs := cmds[cmdName]
	fs.Parse(cmdArgs)

	// Run "scan" command
	if cmdName == "scan" {
		runScan(*dir, *ext, *workers)
	}
}

func runScan(dir, ext string, workers int) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	paths := scanner.Walk(dir, ext)
	results := pool.Start(paths, workers, counter.CountLines)

	totalFiles, totalLines := 0, 0
	for r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "error %s: %v\n", r.Path, r.Err)
			continue
		}
		fmt.Printf("%s: %d\n", r.Path, r.Lines)
		totalFiles++
		totalLines += r.Lines
	}
	fmt.Printf("\nScanned %d files, %d total lines\n", totalFiles, totalLines)
}

func isHelp(s string) bool {
	return s == "-h" || s == "--help" || s == "help"
}

func printUsage() {
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
