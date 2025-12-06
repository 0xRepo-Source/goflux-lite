# GoFlux Lite

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/0xRepo-Source/goflux-lite?style=flat&logo=github)](https://github.com/0xRepo-Source/goflux-lite/releases)
[![License](https://img.shields.io/github/license/0xRepo-Source/goflux-lite?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xRepo-Source/goflux-lite)](https://goreportcard.com/report/github.com/0xRepo-Source/goflux-lite)

**Simple file transfer tools** - Three lightweight binaries for basic file operations.

## Components

**Server** (`gfl-server.exe`) - Minimal file server without web UI  
**Client** (`gfl.exe`) - Command-line file operations  
**Admin** (`gfl-admin.exe`) - Token management  

## Features

**Minimal** - Basic functionality only, no web UI or extras  
**Secure** - Path traversal protection and token authentication  
**Fast** - Lightweight binaries with progress tracking  
**Simple** - Each tool does one thing well  
**Resume** - Resumable uploads for large files  
**Auto-discovery** - Find servers automatically on local network  
**Progress tracking** - Visual progress bars with speed display  
**Auto-firewall** - Automatic Windows Firewall configuration  
**Wildcard support** - Upload multiple files using glob patterns (*, ?, [])  
**Transfer verification** - SHA-256 checksums ensure data integrity  

## Behavior Notes

- `gfl.exe` creates and reads `goflux.json` from the directory where the executable lives. Running commands from other working directories no longer leaves additional config files behind.
- Client commands accept paths that contain spaces without additional quoting (for example `gfl put My Files\archive.zip remote/archive.zip`).

## Quick Start

### 1. Auto-Discovery (Recommended)

The easiest way to get started is using auto-discovery:

```bash
# Start server
gfl-server

# On another machine, discover servers
gfl discover

# Configure client for discovered server
gfl config 192.168.1.100:8080

# Start transferring files
gfl put document.pdf files/document.pdf
gfl get files/document.pdf downloaded.pdf
gfl ls files/
```

### 2. Manual Configuration

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
  discover              Discover GoFlux servers on local network
  config <server:port>  Configure client for discovered server
  get <remote> <local>  Download a file
  put <local> <remote>  Upload file(s) - supports wildcards (*, ?, [])
  ls [path]            List files/directories  

Wildcard examples:
  gfl put *.txt uploads/           # Upload all .txt files
  gfl put report*.pdf archives/    # Upload files matching pattern
  gfl put file[123].log logs/      # Upload file1.log, file2.log, file3.log

By default the client reads `goflux.json` from the directory where `gfl.exe` resides. Use `-config` to point at an alternate configuration file.
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
# Auto-discovery workflow
.\gfl.exe discover
.\gfl.exe config 192.168.1.100:8080

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

# Upload multiple files with wildcards
.\gfl.exe put *.txt uploads/
.\gfl.exe put logs/*.log archives/logs/
.\gfl.exe put report-202[34].pdf reports/

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

- **Path Traversal Protection** - Prevents `../` attacks  
- **Token Authentication** - Secure API access  
- **Permission System** - Granular access controls  
- **SHA-256 Checksums** - Automatic integrity verification for all transfers

## Differences from Full GoFlux

**Removed:**
- Advanced progress bars
- Complex configuration options

**Kept:**
- Core file operations (put/get/ls)
- Token authentication
- Resume functionality for large files
- Multi-chunk uploads with session persistence

## Version

GoFlux Lite v0.1.0 - First stable release