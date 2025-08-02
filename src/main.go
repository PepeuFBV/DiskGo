package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"DiskGo/src/scanner"
	"DiskGo/src/tree"
	"DiskGo/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const maxCPUs = 12
const userHomeDirAsRoot = true

func main() {
   a := app.New()
   w := a.NewWindow("Disk Tree")

   fmt.Println(runtime.NumCPU(), "CPU cores available, using", maxCPUs, "cores for scanning.")
   runtime.GOMAXPROCS(maxCPUs)
   fmt.Println("Starting scan...")

   result := make(chan *tree.Node)

   go func() {
	   var waitgroup sync.WaitGroup // global waitgroup for concurrent goroutines
	   waitgroup.Add(1) // only add 1 for the root, child goroutines add for themselves
	   root, _ := scanner.ScanAllDirectories(GetRootDirectory(userHomeDirAsRoot), 0, &waitgroup)
	   waitgroup.Wait() // wait for all goroutines to finish

	   result <- root
   }()

   root := <-result // receive the result from the channel

   if root == nil {
	   fmt.Println("No directory found.")
	   return
   }

   files, directories := scanner.GetCounters()
   totalBytes := scanner.GetTotalScannedBytes()
   totalNodes := files + directories
   fmt.Printf("Scan completed!\n")
   fmt.Printf("Total nodes in tree: %d\n", totalNodes)
   fmt.Printf("Files found: %d\n", files)
   fmt.Printf("Directories found: %d\n", directories)
   fmt.Printf("Total bytes scanned: %s\n", utils.SizeConverter{Bytes: uint64(totalBytes)}.ToReadable())

   treeWidget := createTreeWidget(root)
   w.SetContent(treeWidget)
   w.ShowAndRun()
}

func GetRootDirectory(useHomeDir bool) string {
   if useHomeDir {
	   return GetUserHomeDirectory()
   }
   if runtime.GOOS == "windows" {
	   return "C:\\"
   }
   return "/"
}

func GetUserHomeDirectory() string {
   homeDir, err := os.UserHomeDir()
   if err != nil {
	   fmt.Println("Error getting user home directory:", err)
	   return ""
   }
   return homeDir
}

func createTreeWidget(root *tree.Node) *widget.Tree {
   nodeMap := make(map[string]*tree.Node)
   buildNodeMap(root, "", nodeMap)

   // determine if an item is a branch (has children)
   isBranch := func(uniqueID string) bool {
	   node, ok := nodeMap[uniqueID]
	   return ok && len(node.Children) > 0
   }

   getChildren := func(uniqueID string) []string {
	   node, ok := nodeMap[uniqueID]
	   if !ok {
		   return []string{}
	   }
	   ids := []string{}
	   for _, child := range node.Children {
		   ids = append(ids, child.ID)
	   }
	   return ids
   }

   createNode := func(branch bool) fyne.CanvasObject {
	   icon := widget.NewIcon(nil)
	   label := widget.NewLabel("")
	   return container.NewHBox(icon, label)
   }

   updateNode := func(uniqueID string, branch bool, obj fyne.CanvasObject) {
	   node := nodeMap[uniqueID]
	   size := utils.SizeConverter{Bytes: uint64(node.Size)}.ToReadable()

	   hbox := obj.(*fyne.Container)
	   icon := hbox.Objects[0].(*widget.Icon)
	   label := hbox.Objects[1].(*widget.Label)

	   if branch {
		   icon.SetResource(theme.FolderIcon()) // closed folder by default
	   } else {
		   icon.SetResource(theme.FileIcon()) // file icon
	   }

	   label.SetText(fmt.Sprintf("%s\t%s", node.Name, size))
   }

   treeWidget := widget.NewTreeWithStrings(map[string][]string{})
   treeWidget.ChildUIDs = getChildren
   treeWidget.IsBranch = isBranch
   treeWidget.CreateNode = createNode
   treeWidget.UpdateNode = updateNode
   treeWidget.Root = root.ID

   return treeWidget
}

// helper function: builds a unique ID for each node based on its path to the root
func buildNodeMap(n *tree.Node, parentID string, m map[string]*tree.Node) {
   if n == nil {
	   return
   }
   if n.ID == "" {
	   n.ID = parentID + "/" + n.Name
   }
   m[n.ID] = n
   for _, child := range n.Children {
	   buildNodeMap(child, n.ID, m)
   }
}

// references:
// - https://docs.fyne.io/
// - https://docs.fyne.io/api/v2.4/widget/tree.html
