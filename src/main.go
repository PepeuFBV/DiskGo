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
	// inicializa o app Fyne e a janela
	a := app.New()
	w := a.NewWindow("Disk Tree")

	// logs de execução do terminal, pasível de remoção após o fim da interface
	fmt.Println(runtime.NumCPU(), "CPU cores disponíveis, usando", maxCPUs, "cores para varredura.")
	runtime.GOMAXPROCS(maxCPUs)
	fmt.Println("Iniciando varredura...")

	resultado := make(chan *tree.Node)

	// Goroutine para chamar a função de busca dos diretórios
	go func() {
		var waitgroup sync.WaitGroup // número de goroutines concorrentes global
		// só adiciona 1 para a raiz, as goroutines filhas adicionam para si mesmas
		waitgroup.Add(1)
		raiz, _ := scanner.BuscarTodosDiretorios(PegarDiretorioRaiz(userHomeDirAsRoot), 0, &waitgroup)
		waitgroup.Wait() // espera todas as goroutines terminarem

		resultado <- raiz
	}()

	// recebe o resultado do canal
	raiz := <- resultado

	if raiz == nil {
		fmt.Println("Nenhum diretório encontrado.")
		return
	}

	arquivos, diretorios := scanner.ObterContadores()
	totalBytes := scanner.ObterTotalBytesVasculhados()
	totalNodes := arquivos + diretorios
	fmt.Printf("Varredura concluída!\n")
	fmt.Printf("Total de nós na árvore: %d\n", totalNodes)
	fmt.Printf("Arquivos encontrados: %d\n", arquivos)
	fmt.Printf("Diretórios encontrados: %d\n", diretorios)
	fmt.Printf("Total de bytes vasculhados: %s\n", utils.SizeConverter{Bytes: uint64(totalBytes)}.ToReadable())

	// instância o Widget Tree e exibe
	treeWidget := criarTreeWidget(raiz) // função estrela da noite
	w.SetContent(treeWidget)
	w.ShowAndRun()
}

// retorna diretório raiz com base no sistema operacional
func PegarDiretorioRaiz(useHomeDir bool) string {
	if useHomeDir {
		return ObterDiretorioHomeUsuario()
	}
	if runtime.GOOS == "windows" {
		return "C:\\"
	}
	return "/"
}

func ObterDiretorioHomeUsuario() string {
	diretorioHome, erro := os.UserHomeDir()
	if erro != nil {
		fmt.Println("Erro ao obter diretório home:", erro)
		return ""
	}
	return diretorioHome
}

// montagem do Widget Tree
func criarTreeWidget(raiz *tree.Node) *widget.Tree {
	// Cria um mapa para armazenar referências a nós por ID
	mapaDeNos := make(map[string]*tree.Node)
	montaMapDeNos(raiz, "", mapaDeNos)

	// aviso: as funções abaixo estão em inglês apenas para combinar com os nomes da documentação

	// função que define se um item é um branch (possui filhos)
	isBranch := func(idUnico string) bool {
		no, ok := mapaDeNos[idUnico]
		return ok && len(no.Children) > 0
	}

	// retorna os filhos de um determinado ID
	getChildren := func(idUnico string) []string {
		no, ok := mapaDeNos[idUnico]
		if !ok {
			return []string{}
		}
		ids := []string{}
		for _, child := range no.Children {
			ids = append(ids, child.ID)
		}
		return ids
	}

	// cria o rótulo visual para cada item da árvore
	createNode := func(branch bool) fyne.CanvasObject {
	icone := widget.NewIcon(nil)
	texto := widget.NewLabel("")
	return container.NewHBox(icone, texto)
}

	updateNode := func(idUnico string, branch bool, obj fyne.CanvasObject) {
		no := mapaDeNos[idUnico]
		tamanho := utils.SizeConverter{Bytes: uint64(no.Size)}.ToReadable()

		hbox := obj.(*fyne.Container)
		icone := hbox.Objects[0].(*widget.Icon)
		texto := hbox.Objects[1].(*widget.Label)

		// define ícone baseado em tipo e branch
		if branch {
			icone.SetResource(theme.FolderIcon()) // pasta "fechada" por padrão
		} else {
			icone.SetResource(theme.FileIcon())
		}

		texto.SetText(fmt.Sprintf("%s	%s", no.Name, tamanho))
	}


	// monta finalmente a árvore
	// popular esses atributos diz à arvore como lidar com a árvore
	treeWidget := widget.NewTreeWithStrings(map[string][]string{})
	treeWidget.ChildUIDs = getChildren
	treeWidget.IsBranch = isBranch
	treeWidget.CreateNode = createNode
	treeWidget.UpdateNode = updateNode
	treeWidget.Root = raiz.ID

	return treeWidget
}

// função auxiliar: monta um ID único para cada nó, baseado no seu caminho até a raiz
func montaMapDeNos(n *tree.Node, parentID string, m map[string]*tree.Node) {
	if n == nil {
		return
	}
	if n.ID == "" {
		n.ID = parentID + "/" + n.Name
	}
	m[n.ID] = n
	for _, child := range n.Children {
		montaMapDeNos(child, n.ID, m)
	}
}

// referências;
//- https://docs.fyne.io/
//- https://docs.fyne.io/api/v2.4/widget/tree.html
