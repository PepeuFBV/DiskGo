package scanner

import (
	"os"
	"path/filepath"
	"sync"

	"DiskGo/src/tree"
)

const maxDepth = -1 // -1 for no limit

func SearchAllDirs(path string, depth int, waitgroup *sync.WaitGroup) (*tree.Node, error) {
    defer waitgroup.Done()

    if maxDepth >= 0 && depth > maxDepth {
        return nil, nil
    }

    info, err := os.Stat(path)
	if err != nil {
        return nil, err
    }

    node := &tree.Node{
        Name: info.Name(),
        Size: 0,
        Type: func() string {
            if info.IsDir() {
                return "directory"
            }
            return "file"
        }(),
    }

	if !info.IsDir() {
		node.Size = info.Size()
		return node, nil
	}

    if info.IsDir() {
        entries, err := os.ReadDir(path)
        if err != nil {
            return nil, err
        }

        var mu sync.Mutex
        var localWG sync.WaitGroup

        for _, entry := range entries {
            childPath := filepath.Join(path, entry.Name())
            localWG.Add(1)
            waitgroup.Add(1)
            go func(p string) {
                defer localWG.Done()
                childNode, err := SearchAllDirs(p, depth+1, waitgroup)
                if err == nil && childNode != nil {
                    mu.Lock()
                    node.Children = append(node.Children, childNode)
                    node.Size += childNode.Size
                    mu.Unlock()
                }
            }(childPath)
        }
        localWG.Wait() // wait for all children to finish before returning
    }

    return node, nil
}
