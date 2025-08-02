package scanner

import (
	"DiskGo/src/tree"
	"DiskGo/src/utils"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

// global atomic counters (mutex-free)
var fileCounter int64
var dirCounter int64
var totalScannedBytes int64
var activeGoroutines int64

const maxDepth = -1 // -1 for unlimited
var maxGoroutines = 300 // limit concurrent goroutines, 0 for unlimited
const maxMemoryBytes = 1024 * 1024 * 1024 * 12 // 12 GB memory limit

const printProgress = true // enable progress log
const printEveryFiles = 100000 // 0 to disable progress log

var semaphore chan struct{} // semaphore to limit concurrent goroutines

func init() { // always initialize the semaphore with a safe value
    if maxGoroutines == 0 {
        maxGoroutines = 200 // safe default value
    }
    semaphore = make(chan struct{}, maxGoroutines)
}

func CheckMemoryLimit() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    if m.Alloc > maxMemoryBytes {
        log.Fatalf("Memory limit exceeded: %d bytes used", m.Alloc)
    }
}

func LogProgress() {
    total := atomic.LoadInt64(&fileCounter)
    if printProgress {
        if printEveryFiles > 0 && total%printEveryFiles == 0 && total != 0 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            memCurrent := utils.SizeConverter{Bytes: m.Alloc}.ToReadable()
            memTotalScanned := utils.SizeConverter{Bytes: uint64(atomic.LoadInt64(&totalScannedBytes))}.ToReadable()
            numGoroutines := atomic.LoadInt64(&activeGoroutines)
            log.Printf("Progress: files=%d, dirs=%d, mem_current=%s, mem_total_scanned=%s, goroutines=%d", total, atomic.LoadInt64(&dirCounter), memCurrent, memTotalScanned, numGoroutines)
        }
    }
}

func ScanAllDirectories(path string, depth int, globalWG *sync.WaitGroup) (*tree.Node, error) {
    defer globalWG.Done()

    CheckMemoryLimit()

    if maxDepth >= 0 && depth > maxDepth {
        return nil, nil
    }

    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil
        }
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
        atomic.AddInt64(&totalScannedBytes, info.Size())
        atomic.AddInt64(&fileCounter, 1)
        LogProgress()
        return node, nil
    } else {
        entries, err := os.ReadDir(path)
        if err != nil {
            return nil, err
        }

        var mutex sync.Mutex
        var localWG sync.WaitGroup

        for _, entry := range entries {
            childPath := filepath.Join(path, entry.Name())
            select {
            case semaphore <- struct{}{}: // acquire semaphore to limit goroutines
                localWG.Add(1)
                globalWG.Add(1)
                atomic.AddInt64(&activeGoroutines, 1)

                go func(p string) {
                    defer localWG.Done()
                    defer func() {
                        atomic.AddInt64(&activeGoroutines, -1)
                        <-semaphore
                        if r := recover(); r != nil {
                            log.Printf("panic in goroutine for %s: %v", p, r)
                        }
                    }()

                    childNode, err := ScanAllDirectories(p, depth+1, globalWG)
                    if err == nil && childNode != nil {
                        mutex.Lock()
                        node.Children = append(node.Children, childNode)
                        node.Size += childNode.Size
                        mutex.Unlock()
                    } else if err != nil && !os.IsNotExist(err) {
                        log.Printf("Error scanning %s: %v", p, err)
                    }
                }(childPath)
            default: // recursively scan without semaphore if limit reached
                globalWG.Add(1)
                childNode, err := ScanAllDirectories(childPath, depth+1, globalWG)
                if err == nil && childNode != nil {
                    mutex.Lock()
                    node.Children = append(node.Children, childNode)
                    node.Size += childNode.Size
                    mutex.Unlock()
                } else if err != nil && !os.IsNotExist(err) {
                    log.Printf("Error scanning %s: %v", childPath, err)
                }
            }
        }
        localWG.Wait()
    }

    if node.Type == "directory" {
        atomic.AddInt64(&dirCounter, 1)
        LogProgress()
    }

    return node, nil
}

func GetCounters() (files int64, directories int64) {
    return atomic.LoadInt64(&fileCounter), atomic.LoadInt64(&dirCounter)
}

func GetTotalScannedBytes() int64 {
    return atomic.LoadInt64(&totalScannedBytes)
}

func GetActiveGoroutines() int64 {
    return atomic.LoadInt64(&activeGoroutines)
}
