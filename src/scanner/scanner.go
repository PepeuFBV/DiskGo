package scanner

import (
	"DiskGo/src/tree"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

var fileCount int64 // atomic counter for files
var dirCount int64  // atomic counter for directories

const maxDepth = -1 // -1 for no limit
const maxGoroutines = 50 // limit concurrent goroutines
const maxMemoryBytes = 1024 * 1024 * 1024 * 12 // 12 GB memory limit

const printEvery = 10000

var semaphore = make(chan struct{}, maxGoroutines)

func checkMemoryLimit() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    if m.Alloc > maxMemoryBytes {
        log.Fatalf("Memory limit exceeded: %d bytes used", m.Alloc)
    }
}

func SearchAllDirs(path string, depth int, waitgroup *sync.WaitGroup) (*tree.Node, error) {
    defer waitgroup.Done() // ensure waitgroup is decremented when function exits

    checkMemoryLimit()

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
                atomic.AddInt64(&dirCount, 1) // increment dir counter
                return "directory"
            }
            atomic.AddInt64(&fileCount, 1) // increment file counter
            return "file"
        }(),
    }

    total := atomic.LoadInt64(&fileCount) + atomic.LoadInt64(&dirCount)
    if total % printEvery == 0 && total != 0 { // log progress every 1000 files/dirs
        log.Printf("Progress: %d files, %d dirs", atomic.LoadInt64(&fileCount), atomic.LoadInt64(&dirCount))
    }

    if !info.IsDir() {
        node.Size = info.Size()
        return node, nil
    } else {
        entries, err := os.ReadDir(path)
        if err != nil {
            return nil, err
        }

        var mutex sync.Mutex // mutex to protect shared state
        var localWG sync.WaitGroup // local waitgroup for goroutines in this function

        for _, entry := range entries {
            childPath := filepath.Join(path, entry.Name())
            select {
            case semaphore <- struct{}{}: // acquire semaphore, and its capacity is not exceeded
                localWG.Add(1)
                waitgroup.Add(1)
                go func(p string) { // goroutine to process child directory
                    defer localWG.Done()
                    defer func() { <-semaphore }() // release semaphore
                    childNode, err := SearchAllDirs(p, depth + 1, waitgroup)
                    if err == nil && childNode != nil {
                        mutex.Lock()
                        node.Children = append(node.Children, childNode)
                        node.Size += childNode.Size
                        mutex.Unlock()
                    }
                }(childPath)
            default: // if semaphore is full, process synchronously in this goroutine
                waitgroup.Add(1)
                childNode, err := SearchAllDirs(childPath, depth + 1, waitgroup)
                if err == nil && childNode != nil {
                    mutex.Lock()
                    node.Children = append(node.Children, childNode)
                    node.Size += childNode.Size
                    mutex.Unlock()
                }
            }
        }
        localWG.Wait() // wait for all goroutines in this function to finish
    }

    return node, nil
}
