package walker_test

import (
	"path/filepath"
	"testing"

	"github.com/jywang99/walker"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
    branches := []string{}
    leaves := []string{}
    wkr := walker.NewWalker(&walker.WalkerConfig{})
    doForLeaf := func(path string, isDir bool) {
        path = filepath.Base(path)
        leaves = append(leaves, path)
    }
    doForDir := func(path string) error {
        path = filepath.Base(path)
        branches = append(branches, path)
        return nil
    }
    wkr.WalkAndDo("./testdata", doForLeaf, doForDir)
    assert.ElementsMatch(t, branches, []string{"testdata"})
    assert.ElementsMatch(t, leaves, []string{"dir1", "dir2", "dir3"})
}

func TestDepth(t *testing.T) {
    branches := []string{}
    leaves := []string{}
    wkr := walker.NewWalker(&walker.WalkerConfig{
        MaxDepth: 1,
        Exts: []string{"txt"},
    })
    doForLeaf := func(path string, isDir bool) {
        path = filepath.Base(path)
        leaves = append(leaves, path)
    }
    doForDir := func(path string) error {
        path = filepath.Base(path)
        branches = append(branches, path)
        return nil
    }
    wkr.WalkAndDo("./testdata", doForLeaf, doForDir)
    assert.ElementsMatch(t, branches, []string{"testdata", "dir1", "dir2", "dir3"})
    assert.ElementsMatch(t, leaves, []string{"file1.txt", "file1_2.txt", "file2_1.txt", "dir3_1"})
}

