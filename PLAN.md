# Go `ls` Reimplementation Plan

This document outlines the roadmap for building a robust, modular, and extensible `ls` clone in Go.

## Milestone 1: The Foundation (Basic Listing)
- [x] **Task 1.1: Project Structure & Setup.** Initialize Go modules and create a package structure that separates CLI argument handling from core file-system logic.
- [x] **Task 1.2: Argument & Flag Parsing.** Implement basic flag parsing (e.g., using the `flag` package) to handle target directories and simple boolean options.
- [x] **Task 1.3: Core Directory Reading.** Use `os.ReadDir` to fetch file entries for a given path and print them to standard output.

## Milestone 2: Data Processing (Filtering & Sorting)
- [x] **Task 2.1: The `FileEntry` Model.** Define a internal struct to store file metadata (Name, Size, Mode, ModTime). This decouples the OS-level data from our display logic.
- [x] **Task 2.2: Hidden Files (`-a`).** Implement a filter to skip files starting with `.` unless the "all" flag is set.
- [x] **Task 2.3: Sorting Engine.** Build a flexible sorting system:
    - Default: Alphabetical (case-insensitive).
    - `-t`: Sort by modification time.
    - `-S`: Sort by file size.
    - `-r`: Reverse the final sorted list.

## Milestone 3: The Long Format (`-l`)
- [x] **Task 3.1: Metadata Expansion.**
    - [x] 3.1.1: Add `Owner`, `Group`, `Nlink`, and `Blocks` fields to `FileEntry`.
    - [x] 3.1.2: Update `NewFileEntry` to extract `syscall.Stat_t` (on Unix) to populate raw IDs and link counts.
- [x] **Task 3.2: Identity Resolution.**
    - [x] 3.2.1: Implement caching for `os/user` lookups to resolve UIDs and GIDs efficiently.
    - [x] 3.2.2: Fallback to numeric IDs if resolution fails.
- [x] **Task 3.3: Permission & Type String.**
    - [x] 3.3.1: Implement a custom formatter to convert `os.FileMode` to the 10-character string (e.g., `drwxr-xr-x`).
- [x] **Task 3.4: Tabular Long Listing.**
    - [x] 3.4.1: Integrate `text/tabwriter` into the `List` function for the `-l` path.
    - [x] 3.4.2: Align columns: Mode, Nlink, Owner, Group, Size, ModTime, Name.
- [x] **Task 3.5: Advanced Formatting.**
    - [x] 3.5.1: Implement `-h` (human-readable) size conversion (B, K, M, G).
    - [x] 3.5.2: Logic for "Recent" vs "Old" time formatting (standard `ls` behavior).

## Milestone 4: Polish & Advanced Features
- [x] **Task 4.1: File Indicators (`-F`).** Implement logic to append type-specific characters: `/` for directories, `*` for executables, `@` for symlinks, etc.
- [x] **Task 4.2: ANSI Color System.** Define a color configuration and apply ANSI escape codes to output based on file type and permissions.
- [x] **Task 4.3: Recursive Listing (`-R`).** Implement a depth-first traversal to list subdirectories, including path headers for each section.
- [x] **Task 4.4: Terminal Detection.** Use `isatty` logic to detect if output is a terminal to enable/disable colors and multi-column mode.
- [x] **Task 4.5: Multi-column Formatting.** Implement a grid-based layout for standard output when not using `-l`, optimizing for terminal width.

## Current Status & Next Steps
### Project Status: Feature Complete (v1.0)
The project has reached its initial goal of becoming a robust, fully-functional `ls` replacement. It features:
- Core `ls` functionality (listing, sorting, hidden files).
- Advanced Long Format (`-l`, `-h`) with proper tabular alignment and Unix metadata.
- Polish features (`-F`, colors, `-R`, TTY-aware multi-column formatting).

## Architectural Principles
1. **Pipe-and-Filter:** Data flows from Scanner -> Filter -> Sorter -> Formatter -> Printer.
2. **Strategy Pattern:** Use interfaces for different sorting and formatting strategies to allow easy addition of new flags.
3. **Platform Independence:** While `ls` is Unix-centric, use Go's standard library abstractions where possible to keep the code clean.

# Build command
Builds must be called `gs`
`mkdir -p builds && go build -ldflags="-s -w" -o builds/gs . && ls -lh builds/gs`
