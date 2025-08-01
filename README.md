

# DiskGo

DiskGo is a Go application designed to make disk/drive reading and analysis easy. The program traverses directories, reads files, and analyzes disk usage in a simple way.

## Features

- Recursive scanning of directories and files
- Displays the directory tree structure with file sizes
- Uses concurrency for faster scanning
- Human-readable size formatting (B, KB, MB, GB, TB)

## Screenshot

![DiskGo Screenshot](images/diskgo-example.png)

## Dependency Installation (Linux)

Before running DiskGo, install the necessary dependencies for Fyne:

```sh
sudo apt-get update
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

## Compilation and Distribution

### Simple Compilation

1. Build the project:

    ```sh
    go build -o diskgo src/main.go
    ```

2. Run the binary:

    ```sh
    ./diskgo
    ```

### Optimized Compilation (Recommended for Distribution)

To create an optimized binary for distribution on any Linux PC:

```sh
# using the build script
./build.sh

# or using Make
make build

# or manually with optimizations
CGO_ENABLED=1 go build -ldflags="-s -w" -o diskgo src/main.go
```

### Multi-Architecture Compilation

To create binaries for different Linux architectures:

```sh
# all supported architectures
make build-all

# or individually:
make build-linux-amd64    # for Intel/AMD 64-bit processors
make build-linux-arm64    # for ARM 64-bit processors
make build-linux-386      # for 32-bit processors
```

### System-wide Installation

To install DiskGo for all system users:

```sh
# using Make (recommended)
make install

# or manually
sudo cp diskgo /usr/local/bin/
```

### Distribution

The compiled binary (`diskgo`) is completely standalone and can be run on any Linux system with the following graphical dependencies installed (Fyne library dependencies):

```sh
# Ubuntu/Debian
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev

# CentOS/RHEL/Fedora
sudo yum install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel

# Arch Linux
sudo pacman -S libx11 libxcursor libxrandr libxinerama libxi mesa
```

## Usage

By default, DiskGo scans the current user's home directory. You can change this in the `userHomeDirAsRoot` variable in the `src/main.go` file (true or false).

### Quick Run (Development)

You can also quickly run the application without compiling using:

```sh
go run src/main.go
```

## Installation

Clone the repository and run using Go:

```sh
git clone https://github.com/pepeufbv/DiskGo.git
cd DiskGo
go run src/main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
