
# DiskGo

DiskGo é um aplicativo em Go projetado para facilitar a leitura e análise de discos/unidades. O programa percorre diretórios, lê arquivos e analisa o uso do disco de forma simples.

## Funcionalidades

- Varredura recursiva de diretórios e arquivos
- Exibe a estrutura de árvore dos diretórios com tamanhos de arquivos
- Usa concorrência para varredura mais rápida
- Formatação de tamanho legível (B, KB, MB, GB, TB)

## Instalação de Dependências (Linux)

Antes de rodar o DiskGo, instale as dependências necessárias para o Fyne:

```sh
sudo apt-get update
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

## Uso

1. Compile o projeto:

    ```sh
    go build -o diskgo src/main.go
    ```

2. Execute o binário:

    ```sh
    ./diskgo
    ```

Por padrão, o DiskGo faz a varredura do diretório home do usuário atual. Você pode mudar isso na variável `userHomeDirAsRoot` no arquivo `src/main.go` (true ou false).

### Execução Rápida

Você também pode rodar rapidamente o aplicativo sem compilar usando:

```sh
go run src/main.go
```

## Instalação

Clone o repositório e rode usando Go:

```sh
git clone https://github.com/pepeufbv/DiskGo.git
cd DiskGo
go run src/main.go
```

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
