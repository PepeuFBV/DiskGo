package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"DiskGo/src/scanner"
	"DiskGo/src/tree"
)

const maxCPUs = 12

const userHomeDirAsRoot = false

func main() {
	fmt.Println(runtime.NumCPU(), "CPU cores available, using", maxCPUs, "cores for scanning.")
	runtime.GOMAXPROCS(maxCPUs)
	fmt.Println("Starting scan...")

	var waitgroup sync.WaitGroup
	result := make(chan *tree.Node)

	waitgroup.Add(1)
	go func() {
		rootNode, _ := scanner.SearchAllDirs(GetRootDir(userHomeDirAsRoot), 0, &waitgroup)
		waitgroup.Wait()
		result <- rootNode // send only when all children are done
	}()

	rootNode := <-result
	if rootNode == nil {
		fmt.Println("No directories found.")
		return
	}

	fmt.Println("Found directories/files:")
	tree.PrintTree(rootNode)
}

func GetRootDir(useHomeDir bool) string {
	if useHomeDir {
		return GetUserHomeDir()
	}

    if runtime.GOOS == "windows" {
        return "C:\\"
    }
    return "/"
}

func GetUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return ""
	}
	return homeDir
}
