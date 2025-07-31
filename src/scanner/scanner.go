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

var contadorArquivos int64 // contador atômico para arquivos
var contadorDiretorios int64  // contador atômico para diretórios
var totalBytesVasculhados int64 // contador atômico para total de bytes vasculhados
var goroutinesAtivas int64 // contador atômico para goroutines ativas do programa

const profundidadeMaxima = -1 // -1 para sem limite
var maximoGoroutines = 150 // limitar goroutines concorrentes, -1 para sem limite
const maximoMemoriaBytes = 1024 * 1024 * 1024 * 12 // limite de memória de 12 GB

const imprimirProgresso = true // habilitar log de progresso
const imprimirCadaArquivos = 100000 // 0 para desabilitar log de progresso

var semaforo chan struct{} // semáforo para limitar goroutines concorrentes

func init() {
    // sempre inicializa o semáforo com um valor seguro
    if maximoGoroutines <= 0 {
        maximoGoroutines = 100 // valor seguro padrão reduzido
    }
    semaforo = make(chan struct{}, maximoGoroutines)
}

func VerificarLimiteMemoria() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    if m.Alloc > maximoMemoriaBytes {
        log.Fatalf("Limite de memória excedido: %d bytes usados", m.Alloc)
    }
}

func LogarProgresso() {
    total := atomic.LoadInt64(&contadorArquivos)
    if imprimirProgresso {
        if total % imprimirCadaArquivos == 0 && total != 0 {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            memAtual := utils.SizeConverter{Bytes: m.Alloc}.ToReadable()
            memTotalVasculhada := utils.SizeConverter{Bytes: uint64(atomic.LoadInt64(&totalBytesVasculhados))}.ToReadable()
            numGoroutines := atomic.LoadInt64(&goroutinesAtivas)
            log.Printf("Progresso: arquivos=%d, dirs=%d, mem_atual=%s, mem_total_vasculhada=%s, goroutines=%d", total, atomic.LoadInt64(&contadorDiretorios), memAtual, memTotalVasculhada, numGoroutines)
        }
    }
}

func BuscarTodosDiretorios(caminho string, profundidade int, grupoEspera *sync.WaitGroup) (*tree.Node, error) {
    defer grupoEspera.Done() // marca a goroutine como concluída

    VerificarLimiteMemoria()

    if profundidadeMaxima >= 0 && profundidade > profundidadeMaxima {
        return nil, nil
    }

    info, err := os.Stat(caminho)
    if err != nil { // se for erro de arquivo inexistente, apenas ignore silenciosamente
        if os.IsNotExist(err) {
            return nil, nil
        }
        return nil, err
    }

    no := &tree.Node{
        Name: info.Name(),
        Size: 0,
        Type: func() string {
            if info.IsDir() {
                atomic.AddInt64(&contadorDiretorios, 1)
                return "diretorio"
            }
            atomic.AddInt64(&contadorArquivos, 1)
            return "arquivo"
        }(),
    }
    LogarProgresso()

    if !info.IsDir() {
        no.Size = info.Size()
        atomic.AddInt64(&totalBytesVasculhados, info.Size())
        return no, nil
    } else {
        entradas, err := os.ReadDir(caminho)
        if err != nil {
            return nil, err
        }

        var mutex sync.Mutex // mutex para proteger o estado compartilhado
        var localWG sync.WaitGroup // waitgroup local para goroutines nesta função

        for _, entrada := range entradas {
            caminhoFilho := filepath.Join(caminho, entrada.Name())
            select {
            case semaforo <- struct{}{}: // se o semáforo não estiver cheio, processa em goroutine
                localWG.Add(1)
                grupoEspera.Add(1)
                atomic.AddInt64(&goroutinesAtivas, 1) // incrementa contador de goroutines ativas
                go func(p string) { // goroutine para processar cada filho
                    defer localWG.Done()
                    defer func() {
                        atomic.AddInt64(&goroutinesAtivas, -1) // decrementa contador de goroutines ativas
                        <-semaforo // libera o semáforo
                        if r := recover(); r != nil {
                            log.Printf("panic in goroutine for %s: %v", p, r)
                        }
                    }()

                    noFilho, err := BuscarTodosDiretorios(p, profundidade+1, grupoEspera)
                    if err == nil && noFilho != nil {
                        mutex.Lock()
                        no.Children = append(no.Children, noFilho)
                        no.Size += noFilho.Size
                        mutex.Unlock()
                    } else if err != nil && !os.IsNotExist(err) {
                        log.Printf("Erro ao buscar %s: %v", p, err)
                    }
                }(caminhoFilho)
            default: // se o semáforo estiver cheio, processa sincronamente nesta goroutine
                grupoEspera.Add(1)
                noFilho, err := BuscarTodosDiretorios(caminhoFilho, profundidade+1, grupoEspera)
                if err == nil && noFilho != nil {
                    mutex.Lock()
                    no.Children = append(no.Children, noFilho)
                    no.Size += noFilho.Size
                    mutex.Unlock()
                } else if err != nil && !os.IsNotExist(err) {
                    log.Printf("Erro ao buscar %s: %v", caminhoFilho, err)
                }
            }
        }
        localWG.Wait() // espera todas as goroutines nesta função terminarem
    }

    return no, nil
}

func ObterContadores() (arquivos int64, diretorios int64) {
    return atomic.LoadInt64(&contadorArquivos), atomic.LoadInt64(&contadorDiretorios)
}

func ObterTotalBytesVasculhados() int64 {
    return atomic.LoadInt64(&totalBytesVasculhados)
}

func ObterGoroutinesAtivas() int64 {
    return atomic.LoadInt64(&goroutinesAtivas)
}
