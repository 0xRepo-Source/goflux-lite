# GoFlux Lite

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/0xRepo-Source/goflux-lite?style=flat&logo=github)](https://github.com/0xRepo-Source/goflux-lite/releases)
[![License](https://img.shields.io/github/license/0xRepo-Source/goflux-lite?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xRepo-Source/goflux-lite)](https://goreportcard.com/report/github.com/0xRepo-Source/goflux-lite)

**Simple file transfer tools** - Three lightweight binaries for basic file operations.

## Components

🖥️ **Server** (`gfl-server.exe`) - Minimal file server without web UI  
📁 **Client** (`gfl.exe`) - Command-line file operations  
🔑 **Admin** (`gfl-admin.exe`) - Token management  

## Features

✨ **Minimal** - Basic functionality only, no web UI or extras  
🔐 **Secure** - Same security fixes as GoFlux v0.4.2  
🚀 **Fast** - Lightweight binaries  
📦 **Simple** - Each tool does one thing well  
🔄 **Resume** - Resumable uploads for large files  

## Quick Start

### 1. Build
```powershell
# Windows
.\build.ps1

# Linux/Mac
chmod +x build.sh
./build.sh
```

### 2. Start Server
```bash
.\gfl-server.exe -port 8080
# or
.\gfl-server.exe -config goflux.json
```

### 3. Create Token
```bash
.\gfl-admin.exe create -user admin -permissions * -days 365
# Copy the token output
```

### 4. Configure Client
Edit `goflux.json` and add your token:
```json
{
  "client": {
    "server_url": "http://localhost:8080",
    "token": "your-token-here"
  }
}
```

**Alternative:** Set environment variable (takes precedence over config file):
```bash
# Windows
$env:GOFLUX_TOKEN_LITE = "your-token-here"

# Linux/Mac
export GOFLUX_TOKEN_LITE="your-token-here"
```

### 5. Use Client
```bash
# Upload file
.\gfl.exe put document.pdf files/document.pdf

# Download file  
.\gfl.exe get files/document.pdf downloaded.pdf

# List files
.\gfl.exe ls files/
```

## Commands

### Server (`gfl-server.exe`)
```bash
gfl-server.exe [-config goflux.json] [-port 8080] [-version]
```

### Client (`gfl.exe`)
```bash
gfl.exe [-config goflux.json] <command> [args...]

Commands:
  get <remote> <local>  Download a file
  put <local> <remote>  Upload a file
  ls [path]            List files/directories  

Note: rm (remove) and mkdir not available in lite version
```

### Admin (`gfl-admin.exe`)
```bash
gfl-admin.exe <command> [options]

Commands:
  create -user <name> [-permissions <perms>] [-days <days>] [-file <tokens.json>]
  list [-file <tokens.json>]
  revoke <token_id> [-file <tokens.json>]
```

## Configuration

Create `goflux.json`:
```json
{
  "server": {
    "address": "localhost:8080",
    "storage_dir": "./data",
    "meta_dir": "./.goflux-meta", 
    "tokens_file": "tokens.json"
  },
  "client": {
    "server_url": "http://localhost:8080",
    "chunk_size": 1048576,
    "token": ""
  }
}
```

## Examples

```bash
# Start server on custom port
.\gfl-server.exe -port 9000

# Create admin token  
.\gfl-admin.exe create -user alice -permissions * -days 90

# Create limited user
.\gfl-admin.exe create -user bob -permissions upload,download -days 30

# List all tokens
.\gfl-admin.exe list

# Upload with progress
.\gfl.exe put largefile.zip backups/largefile.zip

# Upload using environment token (no config needed)
$env:GOFLUX_TOKEN_LITE = "your-token-here"
.\gfl.exe put document.pdf files/document.pdf

# Download to current directory
.\gfl.exe get backups/largefile.zip ./largefile.zip

# Browse directories
.\gfl.exe ls
.\gfl.exe ls backups/
.\gfl.exe ls backups/2024/
```

## Security

- ✅ **Path Traversal Protection** - Prevents `../` attacks  
- ✅ **Token Authentication** - Secure API access
- ✅ **Permission System** - Granular access controls

## Differences from Full GoFlux

**Removed:**
- Advanced progress bars
- Complex configuration options

**Kept:**
- Core file operations (put/get/ls)
- Security fixes from v0.4.2
- Token authentication  
- Resume functionality for large files
- Multi-chunk uploads with session persistence

## Version

GoFlux Lite v0.4.2 - Based on GoFlux v0.4.2 security release