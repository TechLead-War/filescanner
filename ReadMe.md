# FileScanner

`FileScanner` is a modular CLI tool built in Go to scan directories recursively and count lines in files with a specific extension. It is designed with concurrency and extensibility in mind.

---

## Features

- Recursive directory scanning
- File filtering by extension
- Line counting using efficient buffered IO
- Parallel processing via configurable goroutines
- Subcommand architecture for future extensibility

---

## Build Instructions

```bash

git clone git@github.com:TechLead-War/filescanner.git
cd filescanner
go build -o filescanner main.go

Once compied, you can run the binary directly:
  ./filescanner
```

---

## Usage

```bash

./filescanner [command] [flags]
```

If no command is specified, it defaults to `scan`.

### Example

```bash

./filescanner scan -dir=sample -ext=.txt -workers=4
```

---

## Commands

### scan

Recursively scans the given directory and counts the number of lines in each file matching the specified extension.

#### Syntax

```bash

./filescanner scan [flags]
```

#### Flags

| Flag       | Type   | Default | Description                                           |
|------------|--------|---------|-------------------------------------------------------|
| `-dir`     | string | "."     | Directory to recursively scan                        |
| `-ext`     | string | `.txt`  | File extension to match (e.g., `.log`, `.md`)        |
| `-workers` | int    | `0`     | Number of goroutines to use (0 = number of CPU cores) |

#### Examples

```bash

# Basic scan with defaults (current directory, .txt, CPU goroutines)
./filescanner scan

# Scan a log directory using 8 workers
./filescanner scan -dir=/var/log -ext=.log -workers=8

# Using CLI without explicitly specifying command
./filescanner -dir=data -ext=.md -workers=2
```

---

## Help

```bash

./filescanner --help
```

Displays all available commands and flags.

```bash

./filescanner scan --help
```

Displays detailed flag usage for the `scan` command.

---

## Output Format

Each line contains the file path and the line count:

```
sample/a.txt: 3
sample/sub/b.txt: 12

Scanned 2 files, 15 total lines
```

Errors (e.g., unreadable files) are printed to `stderr` but do not stop the execution.

---

## Future scope:
1. Implement COBRA CLI tools for better command line interface.
2. 