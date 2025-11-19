# gfl - File Transfer Client

The `gfl` client provides command-line access to GoFlux Lite server for uploading, downloading, and listing files. It supports resumable uploads, authentication, auto-discovery, and progress tracking.

## Quick Start

### Auto-Discovery (Recommended)
```bash
# Discover servers on network
.\gfl.exe discover

# Configure for discovered server  
.\gfl.exe config 192.168.1.100:8080

# Start transferring files
.\gfl.exe put document.pdf files/document.pdf
.\gfl.exe get files/document.pdf ./document.pdf
.\gfl.exe ls
```

### Manual Configuration
```bash
# Upload a file
.\gfl.exe put document.pdf files/document.pdf

# Download a file
.\gfl.exe get files/document.pdf ./document.pdf

# List files
.\gfl.exe ls
.\gfl.exe ls files/
```

## Commands

### discover - Find Servers
Scans the local network for GoFlux Lite servers.

**Syntax:**
```bash
gfl discover
```

**Example:**
```bash
.\gfl.exe discover
```

**Output:**
```
Discovering GoFlux servers on local network...

Found 2 GoFlux server(s):

1. GoFlux Lite Server (v0.1.0-lite)
   Address: 192.168.1.100:8080
   Status:  Auth Required
   Seen:    just now

2. GoFlux Lite Server (v0.1.0-lite)  
   Address: 192.168.1.50:9000
   Status:  No Auth
   Seen:    5s ago

Use 'gfl config <address>' to configure your client for a server.
```

### config - Auto Configuration  
Automatically configures the client for a discovered server.

**Syntax:**
```bash
gfl config <server_address>
```

**Example:**
```bash
.\gfl.exe config 192.168.1.100:8080
```

**Output:**
```
Configuring client for server: 192.168.1.100:8080
✓ Configuration saved to goflux.json

⚠️ This server requires authentication.
   Set GOFLUX_TOKEN_LITE environment variable or edit goflux.json
   Contact the server administrator for a token.
```

### put - Upload Files
Uploads a local file to the server.

**Syntax:**
```bash
gfl put <local_file> <remote_path> [options]
```

**Options:**
- `-config <path>` - Configuration file (default: "goflux.json")
- `-version` - Show version information

**Examples:**
```bash
# Upload to specific directory
.\gfl.exe put report.pdf reports/2024/report.pdf

# Upload with custom config
.\gfl.exe put data.zip backups/data.zip -config myconfig.json

# Upload large file (automatic chunking and resume)
.\gfl.exe put bigfile.iso downloads/bigfile.iso
```

**Features:**
- **Automatic chunking** for large files
- **Resume support** for interrupted uploads
- **Progress tracking** during transfer
- **Path creation** - creates remote directories as needed

### get - Download Files
Downloads a file from the server to local storage.

**Syntax:**
```bash
gfl get <remote_path> <local_file> [options]
```

**Options:**
- `-config <path>` - Configuration file (default: "goflux.json")
- `-version` - Show version information

**Examples:**
```bash
# Download to specific location
.\gfl.exe get backups/data.zip ./downloads/data.zip

# Download to current directory
.\gfl.exe get files/document.pdf ./document.pdf

# Download with custom config
.\gfl.exe get logs/app.log ./app.log -config myconfig.json
```

**Features:**
- **Streaming download** for efficient memory usage
- **Automatic directory creation** for local paths
- **File integrity** preservation

### ls - List Files
Lists files and directories on the server.

**Syntax:**
```bash
gfl ls [remote_path] [options]
```

**Options:**
- `-config <path>` - Configuration file (default: "goflux.json")
- `-version` - Show version information

**Examples:**
```bash
# List root directory
.\gfl.exe ls

# List specific directory
.\gfl.exe ls backups/
.\gfl.exe ls reports/2024/

# List with custom config
.\gfl.exe ls files/ -config myconfig.json
```

**Output Format:**
- Files and directories listed one per line
- Directories may be indicated by trailing `/` (server dependent)
- Sorted alphabetically

## Authentication

### Configuration File Method
Set your token in the configuration file:

**goflux.json:**
```json
{
  "client": {
    "server_url": "http://192.168.1.100:8080",
    "chunk_size": 1048576,
    "token": "your-token-here"
  }
}
```

### Environment Variable Method (Recommended)
Set the `GOFLUX_TOKEN_LITE` environment variable:

**Windows:**
```powershell
$env:GOFLUX_TOKEN_LITE = "your-token-here"
.\gfl.exe put file.txt remote/file.txt
```

**Linux/Mac:**
```bash
export GOFLUX_TOKEN_LITE="your-token-here"
./gfl put file.txt remote/file.txt
```

**Priority:** Environment variable takes precedence over config file token.

### Getting Tokens
Tokens are created using the `gfl-admin` tool:

```bash
# Create a token with upload/download permissions
.\gfl-admin.exe create -user myuser -permissions upload,download,list -days 30
```

## Configuration

### Complete Configuration Example
```json
{
  "client": {
    "server_url": "http://192.168.1.100:8080",
    "chunk_size": 1048576,
    "token": "optional-fallback-token"
  }
}
```

### Configuration Options

**server_url** - GoFlux Lite server URL
- Format: `http://host:port` or `https://host:port`
- Examples: `"http://localhost:8080"`, `"https://files.company.com"`

**chunk_size** - Upload chunk size in bytes
- Default: `1048576` (1MB)
- Larger chunks = fewer HTTP requests
- Smaller chunks = better resume granularity

**token** - Authentication token (optional)
- Fallback when `GOFLUX_TOKEN_LITE` not set
- Can be empty if server has authentication disabled

## Resumable Uploads

The client automatically handles resumable uploads for large files:

### How It Works
1. **File Chunking** - Large files split into chunks
2. **Upload Tracking** - Server tracks received chunks
3. **Interruption Handling** - If interrupted, upload can resume
4. **Status Check** - Client queries server for missing chunks
5. **Resume Upload** - Only missing chunks are uploaded

### Resume Process
```bash
# Start upload (may be interrupted)
.\gfl.exe put largefile.iso backups/largefile.iso
# ... connection lost ...

# Resume automatically on retry
.\gfl.exe put largefile.iso backups/largefile.iso
# Client detects existing upload and resumes
```

### Progress Indicators
During upload, the client shows:
- Upload progress
- Current chunk being transferred
- Transfer speed
- Estimated time remaining

## Error Handling

### Common Error Messages

**"Authentication required"**
- Server requires authentication
- Set `GOFLUX_TOKEN_LITE` or configure token in config file

**"Permission denied"**
- Token lacks required permissions
- Check token permissions with `gfl-admin list`
- Create new token with correct permissions

**"File not found"**
- Remote file doesn't exist (for download)
- Check path and spelling
- Use `gfl ls` to verify file location

**"Connection refused"**
- Server is not running
- Check server URL and port
- Verify network connectivity

### Debugging Tips
1. **Check server status** - Ensure server is running
2. **Verify authentication** - Test token with simple `ls` command
3. **Check permissions** - Ensure token has required permissions for operation
4. **Test connectivity** - Try accessing server URL in web browser
5. **Review paths** - Use forward slashes, check for typos

## Advanced Usage

### Automation and Scripting
```bash
# Set token once for session
$env:GOFLUX_TOKEN_LITE = "your-token-here"

# Batch operations
.\gfl.exe put report1.pdf reports/report1.pdf
.\gfl.exe put report2.pdf reports/report2.pdf
.\gfl.exe put report3.pdf reports/report3.pdf

# Backup script example
foreach ($file in Get-ChildItem "*.log") {
    .\gfl.exe put $file.Name "logs/$($file.Name)"
}
```

### Integration with CI/CD
```yaml
# GitHub Actions example
- name: Upload build artifacts
  env:
    GOFLUX_TOKEN_LITE: ${{ secrets.GOFLUX_TOKEN }}
  run: |
    ./gfl put dist/app.exe releases/v1.0.0/app.exe
    ./gfl put dist/app.exe.sha256 releases/v1.0.0/app.exe.sha256
```

### Backup Workflows
```bash
# Daily backup script
$date = Get-Date -Format "yyyy-MM-dd"
$env:GOFLUX_TOKEN_LITE = $env:BACKUP_TOKEN

# Upload database backup
.\gfl.exe put "backup-$date.sql" "backups/database/backup-$date.sql"

# Upload file archives
.\gfl.exe put "files-$date.zip" "backups/files/files-$date.zip"
```

## Performance Tips

### Optimal Chunk Size
- **Small files (< 10MB):** Use default 1MB chunks
- **Large files (> 100MB):** Consider 2-4MB chunks
- **Slow networks:** Use smaller chunks (512KB) for better resume
- **Fast networks:** Use larger chunks (4-8MB) for efficiency

### Network Considerations
- **Unstable connections:** Smaller chunks allow better resume
- **High latency:** Larger chunks reduce round trips
- **Bandwidth limits:** Adjust chunk size to fit transfer windows

### Concurrent Operations
```bash
# Multiple uploads in parallel (PowerShell)
$jobs = @()
$jobs += Start-Job { .\gfl.exe put file1.zip backups/file1.zip }
$jobs += Start-Job { .\gfl.exe put file2.zip backups/file2.zip }
$jobs | Wait-Job
```

## Troubleshooting

### Upload Issues
**Upload stalls or fails:**
1. Check available disk space on server
2. Verify upload permissions in token
3. Test with smaller file first
4. Check server logs for errors

**Resume not working:**
1. Ensure same remote path is used
2. Verify server metadata directory is persistent
3. Check that file hasn't been manually deleted on server

### Download Issues
**Download fails:**
1. Verify file exists with `gfl ls`
2. Check download permissions in token
3. Ensure local directory is writable
4. Test network connectivity

### Authentication Issues
**Token not working:**
1. Check token hasn't expired with `gfl-admin list`
2. Verify token hasn't been revoked
3. Ensure correct permissions for operation
4. Test with simple `ls` command first

### Configuration Issues
**Config file not found:**
1. Run from correct directory
2. Use `-config` option to specify path
3. Create config file from example

**Invalid server URL:**
1. Check server is running on specified port
2. Verify protocol (http/https)
3. Test with browser or curl

## Examples by Use Case

### Development Workflow
```bash
# Upload code archives
.\gfl.exe put project-v1.2.3.zip releases/project-v1.2.3.zip

# Download shared libraries
.\gfl.exe get shared/common.lib ./lib/common.lib

# Check available builds
.\gfl.exe ls releases/
```

### Content Management
```bash
# Upload media files
.\gfl.exe put presentation.pptx content/presentations/presentation.pptx

# Download templates
.\gfl.exe get templates/report-template.docx ./templates/

# Browse content library
.\gfl.exe ls content/
.\gfl.exe ls content/images/
```

### System Administration
```bash
# Upload logs for analysis
.\gfl.exe put system.log logs/$(hostname)/system-$(Get-Date -Format yyyy-MM-dd).log

# Download configuration files
.\gfl.exe get config/nginx.conf ./nginx.conf

# Check log directory
.\gfl.exe ls logs/
```
