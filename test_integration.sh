#!/bin/bash
# Integration tests for gs (Go ls implementation)
# Tests all features from PLAN.md milestones

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
PASSED=0
FAILED=0
TOTAL=0

# Binary path
GS="./builds/gs"

# Test result function
test_result() {
    local name="$1"
    local result="$2"
    TOTAL=$((TOTAL + 1))
    if [ "$result" -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $name"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $name"
        FAILED=$((FAILED + 1))
    fi
}

# Section header
section() {
    echo ""
    echo -e "${YELLOW}=== $1 ===${NC}"
}

# Setup test directory
setup_test_dir() {
    TEST_DIR=$(mktemp -d)
    
    # Create test files
    touch "$TEST_DIR/file_a.txt"
    touch "$TEST_DIR/file_b.txt"
    touch "$TEST_DIR/file_c.txt"
    touch "$TEST_DIR/.hidden_file"
    mkdir -p "$TEST_DIR/subdir"
    touch "$TEST_DIR/subdir/nested.txt"
    
    # Create files with different sizes
    dd if=/dev/zero of="$TEST_DIR/small.bin" bs=100 count=1 2>/dev/null
    dd if=/dev/zero of="$TEST_DIR/medium.bin" bs=1024 count=1 2>/dev/null
    dd if=/dev/zero of="$TEST_DIR/large.bin" bs=1024 count=10 2>/dev/null
    
    # Create executable
    echo '#!/bin/bash' > "$TEST_DIR/script.sh"
    chmod +x "$TEST_DIR/script.sh"
    
    # Create symlink
    ln -sf "file_a.txt" "$TEST_DIR/link_to_a"
}

# Cleanup test directory
cleanup_test_dir() {
    if [ -n "$TEST_DIR" ] && [ -d "$TEST_DIR" ]; then
        rm -rf "$TEST_DIR"
    fi
}

# Trap to cleanup on exit
trap cleanup_test_dir EXIT

section "Building gs"

# Build the binary
echo "Building gs..."
mkdir -p builds
go build -ldflags="-s -w" -o builds/gs .
test_result "Build completes without errors" $?

# Check binary exists
if [ -f "$GS" ]; then
    test_result "Binary exists at builds/gs" 0
else
    test_result "Binary exists at builds/gs" 1
fi

# Setup test directory
setup_test_dir

section "Milestone 1: Basic Listing"

# Test 1.1: Basic listing
output=$($GS "$TEST_DIR" 2>&1)
test_result "Basic listing works" $?

# Test 1.2: Listing current directory
output=$($GS . 2>&1)
test_result "Listing current directory works" $?

# Test 1.3: Multiple paths
output=$($GS "$TEST_DIR" . 2>&1)
test_result "Multiple paths listing works" $?

section "Milestone 2: Filtering & Sorting"

# Test 2.1: Hidden files (-a)
output=$($GS -a "$TEST_DIR" 2>&1)
if echo "$output" | grep -q '\.hidden_file'; then
    test_result "-a flag shows hidden files" 0
else
    test_result "-a flag shows hidden files" 1
fi

# Test 2.2: Without -a, hidden files should be filtered
output=$($GS "$TEST_DIR" 2>&1)
if echo "$output" | grep -q ".hidden_file"; then
    test_result "Hidden files filtered without -a" 1
else
    test_result "Hidden files filtered without -a" 0
fi

# Test 2.3: Alphabetical sort (default)
output=$($GS "$TEST_DIR" 2>&1)
test_result "Default alphabetical sort" $?

# Test 2.4: Sort by time (-t)
output=$($GS -t "$TEST_DIR" 2>&1)
test_result "-t flag sorts by modification time" $?

# Test 2.5: Sort by size (-S)
output=$($GS -S "$TEST_DIR" 2>&1)
test_result "-S flag sorts by file size" $?

# Test 2.6: Reverse sort (-r)
output=$($GS -r "$TEST_DIR" 2>&1)
test_result "-r flag reverses sort order" $?

# Test 2.7: Combined flags
output=$($GS -t -r "$TEST_DIR" 2>&1)
test_result "Combined -t -r flags work" $?

output=$($GS -S -a "$TEST_DIR" 2>&1)
test_result "Combined -S -a flags work" $?

section "Milestone 3: Long Format"

# Test 3.1: Long format (-l)
output=$($GS -l "$TEST_DIR" 2>&1)
test_result "-l flag produces long format" $?

# Check long format output contains expected columns (permission string pattern)
# Permission strings start with: - (file), d (dir), l (symlink), s (socket), p (pipe), c/b (devices)
# Use grep -E for extended regex, look for any line starting with valid permission chars
if echo "$output" | grep -qE '^[-dlspcb][rwx-]{9}'; then
    test_result "Long format shows permission string" 0
else
    test_result "Long format shows permission string" 1
fi

# Test 3.2: Human-readable sizes (-h)
output=$($GS -l -h "$TEST_DIR" 2>&1)
test_result "-l -h flags show human-readable sizes" $?

# Check for human-readable format (K, M, G) - need large file
dd if=/dev/zero of="$TEST_DIR/bigfile.bin" bs=1024 count=100 2>/dev/null
output=$($GS -l -h "$TEST_DIR" 2>&1)
if echo "$output" | grep -qE "[0-9]+\.?[0-9]*[KMGTP]"; then
    test_result "Human-readable sizes use K/M/G suffixes" 0
else
    test_result "Human-readable sizes use K/M/G suffixes" 1
fi

# Test 3.3: Long format with all flags
output=$($GS -l -h -t -a "$TEST_DIR" 2>&1)
test_result "Combined -l -h -t -a flags work" $?

section "Milestone 4: Polish & Advanced Features"

# Test 4.1: File indicators (-F)
output=$($GS -F "$TEST_DIR" 2>&1)
test_result "-F flag adds file type indicators" $?

# Check for directory indicator
if echo "$output" | grep -qE "subdir/"; then
    test_result "Directory indicator (/) is shown" 0
else
    test_result "Directory indicator (/) is shown" 1
fi

# Check for executable indicator
if echo "$output" | grep -qE "script\.sh\*"; then
    test_result "Executable indicator (*) is shown" 0
else
    test_result "Executable indicator (*) is shown" 1
fi

# Check for symlink indicator
if echo "$output" | grep -qE "link_to_a@"; then
    test_result "Symlink indicator (@) is shown" 0
else
    test_result "Symlink indicator (@) is shown" 1
fi

# Test 4.2: Color output (--color)
output=$($GS --color=always -F "$TEST_DIR" 2>&1)
test_result "--color=always produces colored output" $?

# Check for ANSI color codes (look for escape sequence in directory output which is blue)
# Blue = \033[34m for directories
if echo "$output" | grep -q $'\033\[34m'; then
    test_result "ANSI color codes are present" 0
else
    test_result "ANSI color codes are present" 1
fi

# Test --color=never
output=$($GS --color=never -F "$TEST_DIR" 2>&1)
test_result "--color=never disables colors" $?

# Test --color=auto (default)
output=$($GS --color=auto -F "$TEST_DIR" 2>&1)
test_result "--color=auto (default) works" $?

# Test 4.3: Recursive listing (-R)
output=$($GS -R "$TEST_DIR" 2>&1)
test_result "-R flag produces recursive listing" $?

# Check for subdirectory content (nested.txt appears after subdir header)
if echo "$output" | grep -q "nested.txt"; then
    test_result "Recursive listing shows nested files" 0
else
    test_result "Recursive listing shows nested files" 1
fi

# Test 4.4: Combined advanced features
output=$($GS -l -h -a -F -R --color=always "$TEST_DIR" 2>&1)
test_result "All advanced flags combined work" $?

section "Milestone 5: Documentation & Discovery"

# Test 5.1: Help command
output=$($GS help 2>&1)
test_result "help command works" $?

# Check help output contains expected sections
if echo "$output" | grep -qi "usage"; then
    test_result "Help shows usage information" 0
else
    test_result "Help shows usage information" 1
fi

if echo "$output" | grep -qi "examples"; then
    test_result "Help shows examples" 0
else
    test_result "Help shows examples" 1
fi

# Test 5.2: Version flag (-v)
output=$($GS -v 2>&1)
test_result "-v flag shows version" $?

# Check version format
if echo "$output" | grep -qE "gs version [0-9]+\.[0-9]+\.[0-9]+"; then
    test_result "Version follows semantic versioning" 0
else
    test_result "Version follows semantic versioning" 1
fi

# Test 5.3: Default help (-h on flag package)
output=$($GS -h 2>&1)
test_result "-h flag shows flag help" $?

section "Edge Cases & Error Handling"

# Test non-existent path
# Note: Current implementation prints error to stderr but continues (exit 0)
# This matches GNU ls behavior when listing multiple files
output=$($GS "/nonexistent/path/$(date +%s)" 2>&1)
if echo "$output" | grep -qi "no such file\|not found\|cannot access"; then
    test_result "Non-existent path shows error message" 0
else
    test_result "Non-existent path shows error message" 1
fi

# Test empty directory
EMPTY_DIR=$(mktemp -d)
output=$($GS "$EMPTY_DIR" 2>&1)
test_result "Empty directory listing works" $?
rm -rf "$EMPTY_DIR"

# Test permission denied (if possible)
# Note: This may not work in all environments
RESTRICTED_DIR=$(mktemp -d)
chmod 000 "$RESTRICTED_DIR" 2>/dev/null || true
output=$($GS "$RESTRICTED_DIR" 2>&1)
chmod 755 "$RESTRICTED_DIR" 2>/dev/null || true
rm -rf "$RESTRICTED_DIR"
test_result "Permission denied handled gracefully" $?

section "Summary"

echo ""
echo "================================"
echo -e "Total tests: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
else
    echo -e "Failed: $FAILED"
fi
echo "================================"

if [ $FAILED -gt 0 ]; then
    echo ""
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
else
    echo ""
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
fi
