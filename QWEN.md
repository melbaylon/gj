# Project Overview

**gj** is a Go reimplementation of the Unix `ls` command. It is a fully-functional `ls` replacement that provides directory listing capabilities with a focus on robustness, modularity, and extensibility.

## Key Features

- **Core functionality**: Directory listing with sorting (alphabetical, by time `-t`, by size `-S`), hidden file handling (`-a`), and reverse order (`-r`)
- **Long format (`-l`)**: Displays file permissions, link count, owner, group, size, modification time, and name in a tabular format
- **Human-readable sizes (`-h`)**: Converts file sizes to K, M, G, etc.
- **File type indicators (`-F`)**: Appends `/` for directories, `*` for executables, `@` for symlinks, `=` for sockets, `|` for pipes
- **ANSI colorization (`--color`)**: Color-coded output based on file type (directories=blue, executables=green, symlinks=cyan, sockets=red, pipes=yellow)
- **Recursive listing (`-R`)**: Lists subdirectories recursively
- **Multi-column output**: Grid-based layout optimized for terminal width when outputting to TTY
- **Built-in help**: Custom `help` command and `--version` flag

## Architecture

The project follows a **pipe-and-filter** architecture with data flowing through: Scanner → Filter → Sorter → Formatter → Printer

### Package Structure

```
gj/
├── main.go              # CLI entry point, flag parsing, version handling
└── internal/ls/
    ├── ls.go            # Core listing logic, formatting, sorting, colorization
    ├── fileentry.go     # FileEntry struct and FormatMode function
    ├── fileentry_unix.go    # Unix-specific metadata extraction (owner/group resolution)
    └── fileentry_windows.go # Windows compatibility layer
```

### Key Components

- **`FileEntry`**: Internal struct that decouples OS-level data from display logic, storing metadata (Name, Size, Mode, ModTime, Owner, Group, Nlink, Blocks)
- **`List()`**: Main entry point in the `ls` package that orchestrates directory scanning, filtering, sorting, and output
- **Platform-specific `NewFileEntry()`**: Uses build tags (`//go:build unix` / `//go:build windows`) for cross-platform compatibility

## Building and Running

### Build Command

```bash
mkdir -p builds && go build -ldflags="-s -w" -o builds/gs . && ls -lh builds/gs
```

The binary is named `gs` (not `gj`).

### Usage

```bash
gs [OPTION]... [FILE]...
```

### Common Commands

```bash
gs -l -h          # List files in long format with human-readable sizes
gs -R             # Recursively list all files and subdirectories
gs --color=always # Force colorized output
gs -t -r          # Sort by modification time, reversed
gs -S -l          # Sort by file size in long format
gs help           # Display help information
gs -v             # Display version information
```

### Available Flags

| Flag | Description |
|------|-------------|
| `-a` | Do not ignore entries starting with `.` |
| `-l` | Use a long listing format |
| `-t` | Sort by modification time |
| `-S` | Sort by file size |
| `-r` | Reverse order while sorting |
| `-F` | Append indicator (one of `*/=>@|`) to entries |
| `-h` | With `-l`, print sizes in human-readable format |
| `-R` | List subdirectories recursively |
| `-color` | Colorize output: `always`, `auto` (default), or `never` |
| `-v` | Display version information and exit |

## Development

### Requirements

- Go 1.24.0 or later
- Dependencies: `golang.org/x/term`, `golang.org/x/sys`

### Testing Practices

- The project uses Go's standard library extensively
- Platform-specific code is isolated using build tags
- Caching is implemented for user/group lookups to improve performance

### CI/CD

GitHub Actions workflow (`.github/workflows/build.yml`) automatically:
- Builds for macOS (Intel, `darwin/amd64`)
- Builds for Windows (`windows/amd64`)
- Creates releases with binaries attached on every push to `main`

## Project Status

**Feature Complete (v1.0.0)** - All milestones from the original plan have been implemented:
- ✅ Milestone 1: Foundation (Basic Listing)
- ✅ Milestone 2: Data Processing (Filtering & Sorting)
- ✅ Milestone 3: The Long Format
- ✅ Milestone 4: Polish & Advanced Features
- ✅ Milestone 5: Documentation & Discovery

## License

MIT License - see `LICENSE` file for details.
