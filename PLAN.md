# Go `ls` Reimplementation Plan

This document outlines the roadmap for building a robust, modular, and extensible `ls` clone in Go.

## Milestone 1: The Foundation (Basic Listing)
- [x] **Task 1.1: Project Structure & Setup.** Initialize Go modules and create a package structure that separates CLI argument handling from core file-system logic.
- [ ] **Task 1.2: Argument & Flag Parsing.** Implement basic flag parsing (e.g., using the `flag` package) to handle target directories and simple boolean options.
- [ ] **Task 1.3: Core Directory Reading.** Use `os.ReadDir` to fetch file entries for a given path and print them to standard output.

## Milestone 2: Data Processing (Filtering & Sorting)
- [ ] **Task 2.1: The `FileEntry` Model.** Define a internal struct to store file metadata (Name, Size, Mode, ModTime). This decouples the OS-level data from our display logic.
- [ ] **Task 2.2: Hidden Files (`-a`).** Implement a filter to skip files starting with `.` unless the "all" flag is set.
- [ ] **Task 2.3: Sorting Engine.** Build a flexible sorting system:
    - Default: Alphabetical (case-insensitive).
    - `-t`: Sort by modification time.
    - `-S`: Sort by file size.
    - `-r`: Reverse the final sorted list.

## Milestone 3: The Long Format (`-l`)
- [ ] **Task 3.1: Metadata Enrichment.** Use `os.Lstat` and the `syscall` (or `os/user`) package to resolve UIDs and GIDs to actual owner and group names.
- [ ] **Task 3.2: Tabular Formatting.** Use `text/tabwriter` to ensure columns align perfectly regardless of varying filename or size lengths.
- [ ] **Task 3.3: Permission String Parsing.** Convert `os.FileMode` into the standard Unix permission string (e.g., `-rw-r--r--`).

## Milestone 4: Polish & Advanced Features
- [ ] **Task 4.1: Indicators & Colors.**
    - `-F`: Append indicators like `/` for directories and `*` for executables.
    - ANSI Colors: Apply colors to output based on file type (e.g., blue for directories, green for executables).
- [ ] **Task 4.2: Recursion (`-R`).** Implement a recursive walk to list all subdirectories and their contents.
- [ ] **Task 4.3: Terminal Intelligence.** Detect if `stdout` is a terminal or a pipe to toggle between multi-column and single-column output automatically.

## Architectural Principles
1. **Pipe-and-Filter:** Data flows from Scanner -> Filter -> Sorter -> Formatter -> Printer.
2. **Strategy Pattern:** Use interfaces for different sorting and formatting strategies to allow easy addition of new flags.
3. **Platform Independence:** While `ls` is Unix-centric, use Go's standard library abstractions where possible to keep the code clean.