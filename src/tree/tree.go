package tree

import (
	"DiskGo/src/utils"
	"fmt"
)

const profundidadeExibicao = 1 // 0 para sem limite de exibição

type Node struct {
	Name     string
	Children []*Node
	Size     int64
	Type     string
    ID       string //para o widget tree
}

func ImprimirArvore(raiz *Node) {
    ImprimirArvoreRecursivo(raiz, "", 0)
}

func ImprimirArvoreRecursivo(raiz *Node, indent string, profundidade int) {
    if raiz == nil {
        return
    }
    fmt.Printf("%sNome: %s, Tipo: %s, Tamanho: %s bytes\n", indent, raiz.Name, raiz.Type, utils.SizeConverter{Bytes: uint64(raiz.Size)}.ToReadable())
    if profundidadeExibicao == 0 || profundidade < profundidadeExibicao {
        for _, filho := range raiz.Children {
            ImprimirArvoreRecursivo(filho, indent + "  ", profundidade + 1)
        }
    }
}

// ContarTotalNodes conta o número total de nós na árvore
func ContarTotalNodes(raiz *Node) int {
    if raiz == nil {
        return 0
    }
    
    total := 1 // conta o nó atual
    for _, filho := range raiz.Children {
        total += ContarTotalNodes(filho)
    }
    
    return total
}
