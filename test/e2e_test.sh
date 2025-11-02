#!/bin/bash
set -e

# E2E tests for confluence-md

TOOL="./confluence-md"
FAILED=0

echo "=== E2E Tests for confluence-md ==="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

test_passed() {
    echo -e "${GREEN}✓${NC} $1"
}

test_failed() {
    echo -e "${RED}✗${NC} $1"
    FAILED=$((FAILED + 1))
}

test_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Test 1: Search returns results
echo "Test 1: Search for 'proposal' returns results"
OUTPUT=$($TOOL search "proposal" --limit 3 2>&1)
if echo "$OUTPUT" | grep -q "Found.*results"; then
    test_passed "Search returned results"
else
    test_failed "Search did not return expected results"
    echo "$OUTPUT"
fi
echo

# Test 2: Search results show page titles
echo "Test 2: Search results contain page titles"
if echo "$OUTPUT" | grep -q "\[1\]"; then
    test_passed "Search results formatted correctly"
else
    test_failed "Search results not formatted correctly"
fi
echo

# Test 3: --lucky flag fetches content
echo "Test 3: --lucky flag fetches page content"
OUTPUT=$($TOOL search "proposal" --lucky 2>&1)
if echo "$OUTPUT" | grep -q "^#"; then
    test_passed "--lucky returned markdown content"
else
    test_failed "--lucky did not return markdown"
    echo "$OUTPUT"
fi
echo

# Test 4: --lucky with --output writes to file
echo "Test 4: --lucky with --output writes to file"
TMPFILE=$(mktemp)
$TOOL search "proposal" --lucky --output "$TMPFILE" 2>/dev/null
if [ -f "$TMPFILE" ] && [ -s "$TMPFILE" ]; then
    test_passed "Output file created with content"
    rm "$TMPFILE"
else
    test_failed "Output file not created or empty"
fi
echo

# Test 5: --index fetches specific result
echo "Test 5: --index fetches specific search result"
OUTPUT=$($TOOL search "proposal" --index 2 2>&1)
if echo "$OUTPUT" | grep -q "^#"; then
    test_passed "--index returned markdown content"
else
    test_failed "--index did not return markdown"
fi
echo

# Test 6: Search with --space filter
echo "Test 6: Search with --space filter"
OUTPUT=$($TOOL search "proposal" --space TECH --limit 3 2>&1)
if echo "$OUTPUT" | grep -q "Space: TECH"; then
    test_passed "--space filter worked"
else
    test_failed "--space filter did not work"
    echo "$OUTPUT"
fi
echo

# Test 7: Search with --mine filter
echo "Test 7: Search with --mine filter"
OUTPUT=$($TOOL search "proposal" --mine --limit 3 2>&1 || true)
if echo "$OUTPUT" | grep -q "Found.*results\|No results found"; then
    test_passed "--mine filter worked (query executed successfully)"
else
    test_failed "--mine filter did not work"
    echo "$OUTPUT"
fi
echo

# Test 8: Search with no results
echo "Test 8: Search with no results handles gracefully"
OUTPUT=$($TOOL search "zzznonexistentqueryzzzz" 2>&1 || true)
if echo "$OUTPUT" | grep -q "No results found"; then
    test_passed "No results handled correctly"
else
    test_failed "No results not handled correctly"
fi
echo

# Test 9: --debug flag provides debug output
echo "Test 9: --debug flag shows debug information"
OUTPUT=$($TOOL search "proposal" --limit 1 --debug 2>&1)
if echo "$OUTPUT" | grep -q "\[DEBUG\]"; then
    test_passed "--debug flag working"
else
    test_failed "--debug flag not working"
fi
echo

# Test 10: --include-metadata includes page metadata
echo "Test 10: --include-metadata includes page info"
OUTPUT=$($TOOL search "proposal" --lucky --include-metadata 2>&1)
if echo "$OUTPUT" | grep -q "Page ID:"; then
    test_passed "--include-metadata working"
else
    test_failed "--include-metadata not working"
    echo "$OUTPUT" | head -20
fi
echo

# Summary
echo "==================================="
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}$FAILED test(s) failed${NC}"
    exit 1
fi
