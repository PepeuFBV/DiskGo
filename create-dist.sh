#!/bin/bash

# DiskGo Distribution Package Creator
# Creates a distributable package with the binary and necessary files

set -e

PACKAGE_NAME="diskgo-linux"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')
DIST_DIR="dist"
PACKAGE_DIR="${DIST_DIR}/${PACKAGE_NAME}-${VERSION}"

echo "Creating distribution package for DiskGo v${VERSION}..."

# Clean previous distributions
if [ -d "$DIST_DIR" ]; then
    echo "Cleaning previous distributions..."
    rm -rf "$DIST_DIR"
fi

# Create package directory
mkdir -p "$PACKAGE_DIR"

# Build the application
echo "Building optimized binary..."
CGO_ENABLED=1 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o "$PACKAGE_DIR/diskgo" \
    src/main.go

# Copy documentation and license
echo "Copying documentation..."
cp README.md "$PACKAGE_DIR/"
cp LICENSE "$PACKAGE_DIR/"

# Create install script
cat > "$PACKAGE_DIR/install.sh" << 'EOF'
#!/bin/bash

echo "Installing DiskGo..."

# Check if running as root for system-wide install
if [ "$EUID" -eq 0 ]; then
    echo "Installing system-wide to /usr/local/bin/"
    cp diskgo /usr/local/bin/
    chmod +x /usr/local/bin/diskgo
    echo "✅ DiskGo installed system-wide"
    echo "Run with: diskgo"
else
    echo "Installing to user directory ~/.local/bin/"
    mkdir -p ~/.local/bin
    cp diskgo ~/.local/bin/
    chmod +x ~/.local/bin/diskgo
    
    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
        echo "Added ~/.local/bin to PATH in ~/.bashrc"
        echo "Please run: source ~/.bashrc"
    fi
    
    echo "✅ DiskGo installed for current user"
    echo "Run with: diskgo (after sourcing ~/.bashrc)"
fi
EOF

chmod +x "$PACKAGE_DIR/install.sh"

# Create uninstall script
cat > "$PACKAGE_DIR/uninstall.sh" << 'EOF'
#!/bin/bash

echo "Uninstalling DiskGo..."

# Check system-wide installation
if [ -f "/usr/local/bin/diskgo" ]; then
    if [ "$EUID" -eq 0 ]; then
        rm /usr/local/bin/diskgo
        echo "✅ Removed system-wide installation"
    else
        echo "❌ Need root privileges to remove system-wide installation"
        echo "Run: sudo ./uninstall.sh"
        exit 1
    fi
fi

# Check user installation
if [ -f "$HOME/.local/bin/diskgo" ]; then
    rm "$HOME/.local/bin/diskgo"
    echo "✅ Removed user installation"
fi

echo "DiskGo uninstalled successfully"
EOF

chmod +x "$PACKAGE_DIR/uninstall.sh"

# Create dependencies install script
cat > "$PACKAGE_DIR/install-deps.sh" << 'EOF'
#!/bin/bash

echo "Installing DiskGo GUI dependencies..."

# Detect the Linux distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
    VER=$VERSION_ID
elif type lsb_release >/dev/null 2>&1; then
    OS=$(lsb_release -si)
    VER=$(lsb_release -sr)
else
    OS=$(uname -s)
    VER=$(uname -r)
fi

echo "Detected OS: $OS"

case $OS in
    *"Ubuntu"*|*"Debian"*)
        echo "Installing dependencies for Ubuntu/Debian..."
        sudo apt-get update
        sudo apt-get install -y libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
        ;;
    *"CentOS"*|*"Red Hat"*|*"RHEL"*)
        echo "Installing dependencies for CentOS/RHEL..."
        sudo yum install -y libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel
        ;;
    *"Fedora"*)
        echo "Installing dependencies for Fedora..."
        sudo dnf install -y libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel
        ;;
    *"Arch"*)
        echo "Installing dependencies for Arch Linux..."
        sudo pacman -S --needed libx11 libxcursor libxrandr libxinerama libxi mesa
        ;;
    *)
        echo "❌ Unsupported distribution: $OS"
        echo "Please install the GUI dependencies manually:"
        echo "- libx11, libxcursor, libxrandr, libxinerama, libxi, mesa/opengl"
        exit 1
        ;;
esac

echo "✅ Dependencies installed successfully"
EOF

chmod +x "$PACKAGE_DIR/install-deps.sh"

# Create README for the package
cat > "$PACKAGE_DIR/INSTALL.md" << EOF
# DiskGo v${VERSION} - Installation Guide

## Quick Start

1. Install GUI dependencies (if not already installed):
   \`\`\`bash
   ./install-deps.sh
   \`\`\`

2. Install DiskGo:
   \`\`\`bash
   # System-wide (requires sudo)
   sudo ./install.sh
   
   # Or user-only
   ./install.sh
   \`\`\`

3. Run DiskGo:
   \`\`\`bash
   diskgo
   \`\`\`

## Manual Installation

If you prefer to install manually:

1. Copy the \`diskgo\` binary to a directory in your PATH
2. Make it executable: \`chmod +x diskgo\`
3. Run with: \`./diskgo\`

## Uninstalling

To remove DiskGo:
\`\`\`bash
./uninstall.sh
\`\`\`

## Requirements

- Linux distribution with GUI support
- OpenGL and X11 libraries (installed by install-deps.sh)

## Support

For issues and support, please visit: https://github.com/PepeuFBV/DiskGo
EOF

# Create the package archive
echo "Creating package archive..."
cd "$DIST_DIR"
tar -czf "${PACKAGE_NAME}-${VERSION}.tar.gz" "${PACKAGE_NAME}-${VERSION}/"

echo "✅ Distribution package created successfully!"
echo "📦 Package: ${DIST_DIR}/${PACKAGE_NAME}-${VERSION}.tar.gz"
echo "📁 Directory: ${PACKAGE_DIR}/"
echo ""
echo "Package contents:"
ls -la "${PACKAGE_NAME}-${VERSION}/"
