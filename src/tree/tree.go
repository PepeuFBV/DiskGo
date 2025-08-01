package tree

import (
	"DiskGo/src/utils"
	"fmt"
)

const displayDepth = 1 // 0 for unlimited display depth

type Node struct {
   Name     string
   Children []*Node
   Size     int64
   Type     string
   ID       string // for the widget tree
}

func PrintTree(root *Node) {
   printTreeRecursive(root, "", 0)
}

func printTreeRecursive(root *Node, indent string, depth int) {
   if root == nil {
       return
   }
   fmt.Printf("%sName: %s, Type: %s, Size: %s bytes\n", indent, root.Name, root.Type, utils.SizeConverter{Bytes: uint64(root.Size)}.ToReadable())
   if displayDepth == 0 || depth < displayDepth {
       for _, child := range root.Children {
           printTreeRecursive(child, indent+"  ", depth+1)
       }
   }
}
