package walker

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// WalkerConfig is the configuration for the Walker
type WalkerConfig struct {
    // Extensions to include
    Exts []string
    // Paths to ignore
    IgnorePaths []string
    // Process dot files or not
    DotFiles bool
    // Max depth to walk
    MaxDepth int

    Logger struct {
        ERROR *log.Logger
        WARN *log.Logger
        INFO *log.Logger
    }
}

// Walker is a directory tree walker
type Walker struct {
    exts map[string]bool
    dotFiles bool
    maxDepth int
    ignorePaths map[string]bool

    elog *log.Logger
    wlog *log.Logger
    ilog *log.Logger
}

// NewWalker creates a new Walker instance
func NewWalker(cfg *WalkerConfig) *Walker {
    return &Walker{
        exts: arrayToMap(cfg.Exts),
        ignorePaths: arrayToMap(cfg.IgnorePaths),
        dotFiles: cfg.DotFiles,
        maxDepth: cfg.MaxDepth,

        elog: cfg.Logger.ERROR,
        wlog: cfg.Logger.WARN,
        ilog: cfg.Logger.INFO,
    }
}

// LeafProcessor is a function that processes a leaf node
// isDir is true if the node is a directory
type LeafProcessor func(string, bool)
// DirProcessor is a function that processes a directory
type DirProcessor func(string) error

// WalkAndDo walks the directory tree rooted at root, calling process for each leaf node and doForDir for each directory
func (wkr *Walker) WalkAndDo(root string, process LeafProcessor, doForDir DirProcessor) error {
    maxDepth := wkr.maxDepth

    var walk func(string, int)
    walk = func(dir string, depth int) {
        if wkr.ignoreEntry(dir) {
            return
        }
        if depth > maxDepth {
            process(dir, true)
            return 
        }

        // process this dir
        err := doForDir(dir)
        if err != nil {
            return
        }

        files, err := os.ReadDir(dir)
        if err != nil {
            wkr.elog.Printf("Error when reading directory: %v\n", err)
            return
        }

        // directory contents
        for _, file := range files {
            // descend into subdirs
            pth := filepath.Join(dir, file.Name())
            if file.IsDir() {
                walk(pth, depth + 1)
                continue
            }

            // process files
            ext := getExt(file.Name())
            if !wkr.ignoreEntry(file.Name()) && len(ext) > 0 && wkr.exts[ext] {
                process(pth, false)
            }
        }
    }
    walk(root, 0)
    return nil
}

// ignoreEntry returns true if the entry should be ignored
func (wkr *Walker) ignoreEntry(path string) bool {
    // dot files
    if !wkr.dotFiles && strings.HasPrefix(filepath.Base(path), ".") {
        return true
    }
    // ignore paths
    if wkr.ignorePaths[filepath.Base(path)] {
        return true
    }

    return false
}

// getExt returns the extension of the file
func getExt(path string) string {
    ext := filepath.Ext(path)
    if ext == "" {
        return ""
    }
    return strings.ToLower(ext[1:])
}

