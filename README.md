
# DiskGo

DiskGo é um aplicativo em Go projetado para facilitar a leitura e análise de discos/unidades. O programa percorre diretórios, lê arquivos e analisa o uso do disco de forma simples.

## Funcionalidades

- Varredura recursiva de diretórios e arquivos
- Exibe a estrutura de árvore dos diretórios com tamanhos de arquivos
- Usa concorrência para varredura mais rápida
- Formatação de tamanho legível (B, KB, MB, GB, TB)

## Screenshot

![DiskGo Screenshot](images/diskgo-example.png)

## Instalação de Dependências (Linux)

Antes de rodar o DiskGo, instale as dependências necessárias para o Fyne:

```sh
sudo apt-get update
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

## Compilação e Distribuição

### Compilação Simples

1. Compile o projeto:

    ```sh
    go build -o diskgo src/main.go
    ```

2. Execute o binário:

    ```sh
    ./diskgo
    ```

### Compilação Otimizada (Recomendado para Distribuição)

Para criar um binário otimizado para distribuição em qualquer PC Linux:

```sh
# Usando o script de build
./build.sh

# Ou usando Make
make build

# Ou manualmente com otimizações
CGO_ENABLED=1 go build -ldflags="-s -w" -o diskgo src/main.go
```

### Compilação para Múltiplas Arquiteturas

Para criar binários para diferentes arquiteturas Linux:

```sh
# Todas as arquiteturas suportadas
make build-all

# Ou individualmente:
make build-linux-amd64    # Para processadores Intel/AMD 64-bit
make build-linux-arm64    # Para processadores ARM 64-bit
make build-linux-386      # Para processadores 32-bit
```

### Instalação Sistema-wide

Para instalar o DiskGo para todos os usuários do sistema:

```sh
# Usando Make (recomendado)
make install

# Ou manualmente
sudo cp diskgo /usr/local/bin/
```

### Distribuição

O binário compilado (`diskgo`) é completamente autônomo e pode ser executado em qualquer sistema Linux com as seguintes dependências gráficas instaladas (dependências da biblioteca Fyne):

```sh
# Ubuntu/Debian
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev

# CentOS/RHEL/Fedora
sudo yum install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel

# Arch Linux
sudo pacman -S libx11 libxcursor libxrandr libxinerama libxi mesa
```

## Uso

Por padrão, o DiskGo faz a varredura do diretório home do usuário atual. Você pode mudar isso na variável `userHomeDirAsRoot` no arquivo `src/main.go` (true ou false).

### Execução Rápida (Desenvolvimento)

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
