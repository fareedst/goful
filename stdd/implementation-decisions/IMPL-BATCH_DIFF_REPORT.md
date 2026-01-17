# [IMPL:BATCH_DIFF_REPORT] Batch Diff Report Implementation

**Status**: Active (Implemented 2026-01-17)

**Cross-References**: [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]

## Overview

Implementation of the batch diff report CLI command that performs non-interactive directory comparison and outputs a structured YAML report.

## Components

### 1. DiffReport and DiffEntry Structs

```go
// DiffEntry represents a single difference found during batch comparison.
// [IMPL:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]
type DiffEntry struct {
    Name   string `yaml:"name"`   // Filename or dirname (with "/" suffix for dirs)
    Path   string `yaml:"path"`   // Relative path from comparison root
    Reason string `yaml:"reason"` // Why it differs (e.g., "size mismatch", "missing in window 2")
    IsDir  bool   `yaml:"isDir"`  // Whether entry is a directory
}

// DiffReport is the YAML-serializable output of a batch diff search.
// [IMPL:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]
type DiffReport struct {
    Directories               []string    `yaml:"directories"`
    TotalFilesChecked         int         `yaml:"totalFilesChecked"`
    TotalDirectoriesTraversed int         `yaml:"totalDirectoriesTraversed"`
    DurationSeconds           float64     `yaml:"durationSeconds"`
    Differences               []DiffEntry `yaml:"differences"`
}
```

### 2. BatchNavigator

Headless implementation of the `Navigator` interface that loads directories without TUI dependencies:

```go
// BatchNavigator implements Navigator for headless batch comparison.
// [IMPL:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]
type BatchNavigator struct {
    dirs        []*Directory
    initialDirs []string
}

// NewBatchNavigator creates a navigator from directory paths.
func NewBatchNavigator(paths []string) (*BatchNavigator, error)

// GetDirs returns current directories.
func (n *BatchNavigator) GetDirs() []*Directory

// ChdirAll changes to named subdirectory in all directories.
func (n *BatchNavigator) ChdirAll(name string)

// ChdirParentAll navigates to parent in all directories.
func (n *BatchNavigator) ChdirParentAll()

// CurrentPath returns the path of the first directory.
func (n *BatchNavigator) CurrentPath() string

// RebuildComparisonIndex is a no-op for batch mode (no coloring needed).
func (n *BatchNavigator) RebuildComparisonIndex()
```

### 3. RunBatchDiffSearch Function

Main entry point that collects all differences:

```go
// RunBatchDiffSearch performs a complete directory comparison and returns a report.
// [IMPL:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]
func RunBatchDiffSearch(paths []string, progressFn func(path string, filesChecked int)) (*DiffReport, error)
```

### 4. CLI Wiring in main.go

```go
var (
    diffReportFlag = flag.Bool("diff-report", false, "Run batch diff report and exit")
    quietFlag      = flag.Bool("quiet", false, "Suppress progress output to stderr")
)

func main() {
    flag.Parse()
    
    if *diffReportFlag {
        // Skip TUI init, run batch mode
        dirs := flag.Args()
        if len(dirs) < 2 {
            fmt.Fprintln(os.Stderr, "Error: --diff-report requires at least 2 directories")
            os.Exit(1)
        }
        
        // Progress reporter goroutine (unless --quiet)
        // ...
        
        report, err := filer.RunBatchDiffSearch(dirs, progressFn)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error:", err)
            os.Exit(1)
        }
        
        // Output YAML to stdout
        yaml.NewEncoder(os.Stdout).Encode(report)
        
        if len(report.Differences) > 0 {
            os.Exit(2)
        }
        os.Exit(0)
    }
    
    // Normal TUI initialization...
}
```

## Progress Reporting

- Ticker-based goroutine writes to stderr every 2 seconds
- Format: `[progress] Scanning path/to/dir/ (N files checked, M differences found)`
- Suppressed when `--quiet` flag is set
- Uses sync.Mutex to safely read current stats from traversal

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, no differences found |
| 1 | Error (invalid args, directory not found, etc.) |
| 2 | Success, differences found |

## Token Coverage `[PROC:TOKEN_AUDIT]`

- `filer/diffsearch.go`: `[IMPL:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]`
- `main.go`: `[IMPL:BATCH_DIFF_REPORT] [ARCH:BATCH_DIFF_REPORT] [REQ:BATCH_DIFF_REPORT]`
- Tests: `TestBatchNavigator_REQ_BATCH_DIFF_REPORT`, `TestRunBatchDiffSearch_REQ_BATCH_DIFF_REPORT`

## Validation Evidence (2026-01-17)

- `go test ./filer/... -run "BATCH_DIFF_REPORT"` → **11 tests pass** (darwin/arm64, Go 1.24.3)
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1429 token references across 78 files.`

### Tests Implemented
| Test Name | Coverage |
|-----------|----------|
| `TestBatchNavigator_GetDirs_REQ_BATCH_DIFF_REPORT` | Directory loading and GetDirs() |
| `TestBatchNavigator_ChdirAll_REQ_BATCH_DIFF_REPORT` | Subdirectory navigation |
| `TestBatchNavigator_ChdirParentAll_REQ_BATCH_DIFF_REPORT` | Parent navigation |
| `TestBatchNavigator_CurrentPath_REQ_BATCH_DIFF_REPORT` | Path reporting |
| `TestBatchNavigator_InvalidDir_REQ_BATCH_DIFF_REPORT` | Error handling for missing dirs |
| `TestRunBatchDiffSearch_NoDifferences_REQ_BATCH_DIFF_REPORT` | Identical directories |
| `TestRunBatchDiffSearch_MissingFile_REQ_BATCH_DIFF_REPORT` | Missing file detection |
| `TestRunBatchDiffSearch_SizeMismatch_REQ_BATCH_DIFF_REPORT` | Size mismatch detection |
| `TestRunBatchDiffSearch_MissingDir_REQ_BATCH_DIFF_REPORT` | Missing directory detection |
| `TestRunBatchDiffSearch_ProgressCallback_REQ_BATCH_DIFF_REPORT` | Progress callback invocation |
| `TestRunBatchDiffSearch_ThreeWay_REQ_BATCH_DIFF_REPORT` | N-way (3+) comparison |

### Files Modified
- `filer/diffsearch.go` - Added `DiffReport`, `DiffEntry`, `BatchNavigator`, `RunBatchDiffSearch`
- `filer/diffsearch_test.go` - Added 11 unit tests
- `main.go` - Added `--diff-report`, `--quiet` flags and `runBatchDiffReport()` function
