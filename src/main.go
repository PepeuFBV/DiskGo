package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"DiskGo/src/scanner"
	"DiskGo/src/tree"
)

const maxCPUs = 4

func main() {
	fmt.Println(runtime.NumCPU(), "CPU cores available")
	runtime.GOMAXPROCS(maxCPUs)
	fmt.Println("Starting scan...")

	var waitgroup sync.WaitGroup
	result := make(chan *tree.Node)

	waitgroup.Add(1)
	go func() {
		rootNode, _ := scanner.SearchAllDirs(GetRootDir(), 0, &waitgroup)
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

func GetRootDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return "/"
	}
	return filepath.Join(home, "repos")
}
