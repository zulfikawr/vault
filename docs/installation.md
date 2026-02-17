# Installation

This guide covers installing Vault on your system.

## System Requirements

- **Go**: 1.25.6 or later (for building from source)
- **Operating System**: Linux, macOS, or Windows (amd64/arm64)
- **Disk Space**: ~50MB for binary, plus data storage
- **Memory**: Minimum 256MB RAM recommended

## Installation Methods

### Method 1: One-Line Install (Linux/macOS)

The easiest way to install Vault:

```bash
curl -fsSL https://raw.githubusercontent.com/zulfikawr/vault/main/install.sh | bash
```

This script will:
1. Download the latest binary for your platform
2. Make it executable
3. Move it to `/usr/local/bin`

### Method 2: Download Pre-built Binary

1. Visit [GitHub Releases](https://github.com/zulfikawr/vault/releases)
2. Download the binary for your platform:
   - `vault-linux-amd64` - Linux 64-bit
   - `vault-linux-arm64` - Linux ARM 64-bit
   - `vault-darwin-amd64` - macOS Intel
   - `vault-darwin-arm64` - macOS Apple Silicon
   - `vault-windows-amd64.exe` - Windows 64-bit

3. Make executable (Linux/macOS):
   ```bash
   chmod +x vault
   sudo mv vault /usr/local/bin/
   ```

4. Verify installation:
   ```bash
   vault version
   ```

### Method 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/zulfikawr/vault.git
cd vault

# Build the binary
go build -o vault ./cmd/vault

# Install globally (optional)
sudo mv vault /usr/local/bin/
```

### Method 4: Using Go Install

```bash
go install github.com/zulfikawr/vault/cmd/vault@latest
```

The binary will be installed to `$GOPATH/bin/vault`.

## Verify Installation

After installation, verify Vault is working:

```bash
vault version
```

Expected output:
```
Vault version 0.7.0
```

## Uninstall

### Linux/macOS

```bash
sudo rm /usr/local/bin/vault
rm -rf ~/vault_data  # Optional: remove data
```

### Windows

```powershell
# Remove from PATH location
Remove-Item "C:\Program Files\vault\vault.exe"
```

## Next Steps

- [Quick Start Guide](./quickstart.md) - Get up and running in 5 minutes
- [Introduction](./introduction.md) - Learn what Vault can do

## Troubleshooting

### "Command not found"

Make sure Vault is in your PATH:

```bash
# Check if vault is in PATH
which vault

# Add to PATH (Linux/macOS)
export PATH=$PATH:/usr/local/bin
```

### Permission Denied

```bash
# Make executable
chmod +x /path/to/vault

# Or install to a user-writable location
mkdir -p ~/bin
mv vault ~/bin/
export PATH=$PATH:~/bin
```

### Wrong Architecture

Download the correct binary for your system:

```bash
# Check your architecture
uname -m

# x86_64 = amd64
# aarch64 = arm64
```
