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
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

const maxCPUs = 12
const userHomeDirAsRoot = true

//Main------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func main() {
	//Inicializa o app Fyne e a janela
	a := app.New()
	w := a.NewWindow("Disk Tree")

	//Logs de execução do terminal, pasível de remoção após o fim da interface
	fmt.Println(runtime.NumCPU(), "CPU cores available, using", maxCPUs, "cores for scanning.")
	runtime.GOMAXPROCS(maxCPUs)
	fmt.Println("Starting scan...")

	result := make(chan *tree.Node)

	//GoRoutine para chamar a função de busca dos diretórios
	go func() {
		var waitgroup sync.WaitGroup
		waitgroup.Add(1)

		rootNode, _ := scanner.SearchAllDirs(GetRootDir(userHomeDirAsRoot), 0, &waitgroup)
		waitgroup.Wait()

		result <- rootNode
	}()

	//Recebe o resultado do canal
	rootNode := <-result

	if rootNode == nil {
		fmt.Println("No directories found.")
		return
	}

	//Instancia o Widget Tree e exibe
	treeWidget := criarTreeWidget(rootNode) //<-- Função estrela da noite
	w.SetContent(treeWidget)
	w.ShowAndRun()
}

// Retorna diretório raiz com base no sistema operacional
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

//Montagem do Widget Tree
func criarTreeWidget(raiz *tree.Node) *widget.Tree {
	// Cria um mapa para armazenar referências a nós por ID
	mapaDeNos := make(map[string]*tree.Node)
	montaMapDeNos(raiz, "", mapaDeNos)

	//Aviso: as funções abaixo estão em inglês apenas para combinar com os nomes da documentação

	// Função que define se um item é um branch (possui filhos)
	isBranch := func(idUnico string) bool {
		no, ok := mapaDeNos[idUnico]
		return ok && len(no.Children) > 0
	}

	//Retorna os filhos de um determinado ID
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

	//Cria o rótulo visual para cada item da árvore
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

		//Define ícone baseado em tipo e branch
		if branch {
			icone.SetResource(theme.FolderIcon()) //Pasta "fechada" por padrão
		} else {
			icone.SetResource(theme.FileIcon())
		}

		texto.SetText(fmt.Sprintf("%s	%s", no.Name, tamanho))
	}


	//Monta finalmente a árvore
	//Popular esses atributos diz à arvore como lidar com a árvore
	treeWidget := widget.NewTreeWithStrings(map[string][]string{})
	treeWidget.ChildUIDs = getChildren
	treeWidget.IsBranch = isBranch
	treeWidget.CreateNode = createNode
	treeWidget.UpdateNode = updateNode
	treeWidget.Root = raiz.ID

	return treeWidget
}

//Função auxiliar: monta um ID único para cada nó, baseado no seu caminho até a raiz
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

//Referências;
//- https://docs.fyne.io/
//- https://docs.fyne.io/api/v2.4/widget/tree.html