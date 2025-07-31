package tree

import (
	"fmt"
	"DiskGo/src/utils"
)

const displayDepth = 1 // 0 for no display limit

type Node struct {
	Name     string
	Children []*Node
	Size     int64
	Type     string
    ID       string //para o widget tree
}

func PrintTree(root *Node) {
    PrintTreeRecursive(root, "", 0)
}

func PrintTreeRecursive(root *Node, indent string, depth int) {
    if root == nil {
        return
    }
    fmt.Printf("%sName: %s, Type: %s, Size: %s bytes\n", indent, root.Name, root.Type, utils.SizeConverter{Bytes: uint64(root.Size)}.ToReadable())
    if displayDepth == 0 || depth < displayDepth {
        for _, child := range root.Children {
            PrintTreeRecursive(child, indent + "  ", depth + 1)
        }
    }
}
