#!/bin/bash
# Setup demo directories for GIF recordings
# Creates /tmp/demo/{alpha,beta,gamma} with ~6 files each
# Some files missing or different between directories to showcase comparison features

DEMO_ROOT="/tmp/goful-demo"

# Clean up existing demo directories
rm -rf "$DEMO_ROOT"

# Create directories
mkdir -p "$DEMO_ROOT/alpha/subdir/another" \
         "$DEMO_ROOT/beta/subdir/another" \
         "$DEMO_ROOT/gamma/subdir/another"

# Create files in alpha (all 6 files present)
echo "content1" > "$DEMO_ROOT/alpha/file1"
echo "content2" > "$DEMO_ROOT/alpha/file2"
echo "content3" > "$DEMO_ROOT/alpha/file3"
echo "content4" > "$DEMO_ROOT/alpha/file4"
echo "content5" > "$DEMO_ROOT/alpha/file5"
echo "content6" > "$DEMO_ROOT/alpha/file6"

# Create files in beta (missing file6, some with different sizes)
echo "content1" > "$DEMO_ROOT/beta/file1"  # same
echo "content2-different" > "$DEMO_ROOT/beta/file2"  # different size
echo "content3" > "$DEMO_ROOT/beta/file3"  # same
echo "content4-different-longer" > "$DEMO_ROOT/beta/file4"  # different size
echo "content5" > "$DEMO_ROOT/beta/file5"  # same
# file6 missing

# Create files in gamma (missing file5, some with different sizes)
echo "content1" > "$DEMO_ROOT/gamma/file1"  # same
echo "content2" > "$DEMO_ROOT/gamma/file2"  # same
echo "content3-different" > "$DEMO_ROOT/gamma/file3"  # different size
echo "content4" > "$DEMO_ROOT/gamma/file4"  # same
# file5 missing
echo "content6" > "$DEMO_ROOT/gamma/file6"  # same

# Create subdirectory files
RELDIR="$DEMO_ROOT/alpha/subdir"
echo "sub1" > "$RELDIR/subfile1"
echo "sub2" > "$RELDIR/subfile2"
echo "sub3" > "$RELDIR/another/subfile3"
echo "sub4" > "$RELDIR/subfile4"
echo "sub5" > "$RELDIR/subfile5"

RELDIR="$DEMO_ROOT/beta/subdir"
echo "sub1" > "$RELDIR/subfile1"
echo "sub3" > "$RELDIR/another/subfile3"
echo "sub3b" > "$RELDIR/another/subfile3b"
echo "sub2-different" > "$RELDIR/subfile2"  # different size
echo "sub4" > "$RELDIR/subfile4"
echo "sub5" > "$RELDIR/subfile5"

RELDIR="$DEMO_ROOT/gamma/subdir"
echo "sub1" > "$RELDIR/subfile1"
# subfile2 missing in gamma
# subfile3 missing in gamma
echo "sub3b" > "$RELDIR/another/subfile3b"
echo "sub4" > "$RELDIR/subfile4"
echo "sub5" > "$RELDIR/subfile5"

echo "Demo directories created at $DEMO_ROOT"
tree "$DEMO_ROOT"
