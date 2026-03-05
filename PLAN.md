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
- [ ] **Task 3.1: Raw ID Extraction.** Update `FileEntry` and `NewFileEntry` to capture raw UID, GID, and Hard Link count using `syscall.Stat_t`.
- [ ] **Task 3.2: Identity Resolution.** Integrate `os/user` to resolve numeric UIDs and GIDs into actual username and group names.
- [ ] **Task 3.3: Permission String Formatter.** Implement logic to convert `os.FileMode` bits into the standard 10-character Unix string (e.g., `drwxr-xr-x`).
- [ ] **Task 3.4: Basic Tabular Layout.** Use `text/tabwriter` to create the basic `-l` columns: permissions, links, owner, group, size, date, and name.
- [ ] **Task 3.5: Time Formatting.** Standardize the modification time display (e.g., "Jan _2 15:04" or "Jan _2  2006") to match standard `ls` behavior.
- [ ] **Task 3.6: Human-Readable Sizes (`-h`).** Add an optional task to format file sizes into KB, MB, GB, etc., when the `-h` flag is provided.

## Milestone 4: Polish & Advanced Features
- [x] **Task 4.1: File Indicators (`-F`).** Implement logic to append type-specific characters: `/` for directories, `*` for executables, `@` for symlinks, etc.
- [x] **Task 4.2: ANSI Color System.** Define a color configuration and apply ANSI escape codes to output based on file type and permissions.
- [x] **Task 4.3: Recursive Listing (`-R`).** Implement a depth-first traversal to list subdirectories, including path headers for each section.
- [x] **Task 4.4: Terminal Detection.** Use `isatty` logic to detect if output is a terminal to enable/disable colors and multi-column mode.
- [x] **Task 4.5: Multi-column Formatting.** Implement a grid-based layout for standard output when not using `-l`, optimizing for terminal width.

## Architectural Principles
1. **Pipe-and-Filter:** Data flows from Scanner -> Filter -> Sorter -> Formatter -> Printer.
2. **Strategy Pattern:** Use interfaces for different sorting and formatting strategies to allow easy addition of new flags.
3. **Platform Independence:** While `ls` is Unix-centric, use Go's standard library abstractions where possible to keep the code clean.

# Build command
Builds must be called `gs`
`mkdir -p builds && go build -ldflags="-s -w" -o builds/gs . && ls -lh builds/gs`
