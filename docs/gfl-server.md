# gfl-server - File Server

The `gfl-server` is a lightweight file server that provides secure file upload, download, and listing operations through a REST API. It supports resumable uploads, token authentication, chunked file transfers, and automatic network discovery.

## Quick Start

```bash
# Start server with default settings (includes auto-discovery)
.\gfl-server.exe

# Start on custom port
.\gfl-server.exe -port 9000

# Start with configuration file
.\gfl-server.exe -config myconfig.json
```

The server automatically announces itself on the local network (UDP port 8081) so clients can discover it using `gfl discover`.

## Command Line Options

- `-config <path>` - Configuration file path (default: "goflux.json")
- `-port <port>` - Server port, overrides config (uses internal IP)
- `-version` - Print version information

## Configuration

The server uses a JSON configuration file (default: `goflux.json`):

```json
{
  "server": {
    "address": "0.0.0.0:8080",
    "storage_dir": "./data",
    "meta_dir": "./.goflux-meta",
    "tokens_file": "tokens.json",
    "tls_cert": "",
    "tls_key": ""
  }
}
```

### Configuration Options

**address** - Listen address and port
- Format: `host:port` or `ip:port`
- Example: `"0.0.0.0:8080"`, `"192.168.1.100:9000"`
- Use `0.0.0.0` to listen on all interfaces

**storage_dir** - Root directory for file storage
- All uploaded files are stored under this directory
- Must be writable by the server process
- Created automatically if it doesn't exist

**meta_dir** - Directory for upload session metadata
- Stores resume information for interrupted uploads
- Used for chunked upload tracking
- Should be persistent across server restarts

**tokens_file** - Authentication token file (optional)
- Path to JSON file containing access tokens
- Leave empty (`""`) to disable authentication
- Created with `gfl-admin` tool

**tls_cert** / **tls_key** - TLS/SSL configuration (optional)
- Paths to certificate and key files for HTTPS
- Leave empty for HTTP-only operation
- Both required for TLS to work

## API Endpoints

### Authentication
**GET /auth/challenge** - Get authentication challenge (if auth enabled)
- Returns nonce for challenge-response authentication
- No authentication required

### Discovery  
**GET /config** - Get server configuration for auto-discovery
- Returns server configuration JSON for client setup
- No authentication required
- Used by `gfl config` command

### File Operations
**POST /upload** - Upload file chunk
- Content-Type: `application/json`
- Body: Chunk data with metadata
- Supports resumable uploads

**GET /upload/status?path=<file_path>** - Check upload status
- Returns completion status and missing chunks
- Used for resume functionality

**GET /download?path=<file_path>** - Download file
- Returns file content
- Content-Type determined by file extension

**GET /list?path=<directory_path>** - List directory contents
- Returns JSON array of files and directories
- Empty path lists root directory

### Authentication Methods

**Bearer Token:**
```
Authorization: Bearer <token>
```

**Challenge-Response:**
```
Authorization: Challenge <response>;<nonce>;<token_id>
```

## Security Features

### Path Traversal Protection
- All file paths are sanitized
- Prevents `../` directory traversal attacks
- Files are confined to the storage directory

### Token Authentication
- Optional but recommended for production
- Granular permission system (upload/download/list)
- Token expiration and revocation support

### Permission System
- `upload` - Allow file uploads
- `download` - Allow file downloads
- `list` - Allow directory listing
- `*` - All permissions (admin)

### Network Discovery
- **UDP Broadcast Service** - Automatically announces server presence
- **Port:** 8081 (UDP) 
- **Interval:** 30 seconds
- **Format:** JSON with server info (name, version, address, auth status)
- **Usage:** Enables `gfl discover` command to find servers

## Startup Messages

**With Authentication (Green):**
```
Authentication enabled (challenge-response supported)
goflux server listening on 192.168.1.100:8080
```

**Without Authentication (Red Warning):**
```
⚠️ Authentication disabled - all endpoints are public!
It is recommended to enable authentication in production environments.
Please run gfl-admin to create token files and enable auth.
goflux server listening on 192.168.1.100:8080
```

## Resume Functionality

The server supports automatic resume for interrupted uploads:

1. **Chunked Uploads** - Files split into manageable chunks
2. **Session Tracking** - Each upload gets a persistent session
3. **Missing Chunk Detection** - Server tracks which chunks are received
4. **Automatic Recovery** - Clients can query status and resume
5. **Metadata Persistence** - Sessions survive server restarts

### Resume Process
1. Client uploads file chunks
2. Server tracks progress in metadata directory
3. If interrupted, client queries `/upload/status`
4. Server returns list of missing chunks
5. Client resumes by uploading only missing chunks

## Production Deployment

### Basic Setup
```bash
# 1. Create configuration
cp goflux.example.json goflux.json

# 2. Create storage directories
mkdir data
mkdir .goflux-meta

# 3. Set up authentication
.\gfl-admin.exe create -user admin -permissions * -days 365

# 4. Start server
.\gfl-server.exe
```

### Security Hardening
1. **Enable Authentication** - Always use tokens in production
2. **Use HTTPS** - Configure TLS certificates for encrypted transport
3. **Restrict Permissions** - Create users with minimal required permissions
4. **Monitor Access** - Review logs and token usage regularly
5. **Backup Tokens** - Keep secure copies of tokens.json
6. **Firewall** - Restrict network access to required clients only

### File System Permissions
```bash
# Secure the data directory
chmod 750 data/
chmod 640 goflux.json
chmod 600 tokens.json
```

## Monitoring and Logs

The server outputs to standard output/error. For production:

```bash
# Run with logging
.\gfl-server.exe > server.log 2>&1

# Run as background service (Linux)
nohup .\gfl-server > server.log 2>&1 &

# Windows Service (consider using NSSM or similar)
```

## Performance Notes

- **Concurrent Uploads** - Server handles multiple simultaneous transfers
- **Memory Efficient** - Chunked processing keeps memory usage low  
- **Resume Friendly** - Interrupted transfers don't lose progress
- **Fast Restarts** - Server state persists across restarts

## Troubleshooting

**Server won't start:**
- Check if port is already in use
- Verify configuration file syntax
- Ensure storage directories are writable

**Authentication not working:**
- Verify tokens.json exists and is valid
- Check token permissions match operation
- Ensure tokens haven't expired

**Upload failures:**
- Check disk space in storage directory
- Verify client has upload permissions
- Review server logs for detailed errors

**Resume not working:**
- Ensure meta_dir is persistent and writable
- Check that session metadata files aren't deleted
- Verify client is sending correct chunk information

## Integration Examples

### Nginx Reverse Proxy
```nginx
location /api/ {
    proxy_pass http://127.0.0.1:8080/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    client_max_body_size 100M;
}
```

### Systemd Service (Linux)
```ini
[Unit]
Description=GoFlux Lite Server
After=network.target

[Service]
Type=simple
User=goflux
WorkingDirectory=/opt/goflux-lite
ExecStart=/opt/goflux-lite/gfl-server
Restart=always

[Install]
WantedBy=multi-user.target
```
