# DiskGo

DiskGo is a Go-based application designed to facilitate disk/drive reading and analysis. The program traverses directories, reads files, and analyzes disk usage with ease.

## Features

- Recursively scans directories and files
- Displays directory tree structure with file sizes
- Uses concurrency for faster scanning
- Human-readable size formatting (B, KB, MB, GB, TB)

## Usage

1. Build the project:

    ```sh
    go build -o diskgo main.go
    ```

2. Run the executable:

    ```sh
    ./diskgo
    ```

By default, DiskGo scans the `repos` directory inside your home folder.

### Quick Start

You can also quickly run the application without building it by using:

```sh
go run main.go
```

## Installation

Clone the repository and build using Go:

```sh
git clone https://github.com/pepeufbv/DiskGo.git
cd DiskGo
go run main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
