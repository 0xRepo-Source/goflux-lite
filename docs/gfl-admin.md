# gfl-admin - Token Management Tool

The `gfl-admin` tool manages authentication tokens for GoFlux Lite server. It allows you to create, list, and revoke access tokens with specific permissions.

## Quick Start

```bash
# Create your first admin token
.\gfl-admin.exe create -user admin -permissions * -days 365

# List all tokens
.\gfl-admin.exe list

# Revoke a token by ID
.\gfl-admin.exe revoke <token_id>
```

## Commands

### create
Creates a new authentication token.

**Syntax:**
```bash
gfl-admin create -user <username> [options]
```

**Options:**
- `-user <name>` - Username for the token (required)
- `-permissions <perms>` - Comma-separated permissions (default: "upload,download,list")
- `-days <days>` - Token validity in days (default: 30)
- `-file <path>` - Token file path (default: "tokens.json")

**Permission Types:**
- `upload` - Allow file uploads
- `download` - Allow file downloads  
- `list` - Allow directory listing
- `*` - All permissions (admin access)

**Examples:**
```bash
# Create admin token (all permissions)
.\gfl-admin.exe create -user alice -permissions * -days 90

# Create upload-only user
.\gfl-admin.exe create -user uploader -permissions upload -days 7

# Create read-only user
.\gfl-admin.exe create -user reader -permissions download,list -days 30

# Use custom token file
.\gfl-admin.exe create -user bob -permissions * -file ./config/tokens.json
```

### list
Displays all tokens with their status and permissions.

**Syntax:**
```bash
gfl-admin list [options]
```

**Options:**
- `-file <path>` - Token file path (default: "tokens.json")

**Output includes:**
- Token ID
- Username
- Permissions
- Created date
- Expiration date
- Status (active/expired/revoked)

**Example:**
```bash
.\gfl-admin.exe list
```

### revoke
Revokes an active token, preventing further use.

**Syntax:**
```bash
gfl-admin revoke <token_id> [options]
```

**Options:**
- `-file <path>` - Token file path (default: "tokens.json")

**Example:**
```bash
# Revoke token by ID (get ID from list command)
.\gfl-admin.exe revoke abc123def456
```

## Token File

Tokens are stored in JSON format (default: `tokens.json`). This file should be:
- Kept secure and backed up
- Readable by the server process
- Protected with appropriate file permissions

**Sample token file structure:**
```json
{
  "tokens": [
    {
      "id": "abc123",
      "token_hash": "sha256_hash_here",
      "user": "admin",
      "permissions": ["*"],
      "created_at": "2024-01-01T00:00:00Z",
      "expires_at": "2025-01-01T00:00:00Z",
      "revoked": false
    }
  ]
}
```

## Security Best Practices

1. **Strong Tokens**: Tokens are automatically generated with cryptographic randomness
2. **Limited Scope**: Grant minimum required permissions
3. **Time Limits**: Use reasonable expiration periods
4. **Regular Rotation**: Periodically create new tokens and revoke old ones
5. **Secure Storage**: Protect the tokens.json file with proper permissions
6. **Monitoring**: Regularly review active tokens with `list` command

## Common Workflows

### Initial Setup
```bash
# 1. Create admin token for server management
.\gfl-admin.exe create -user admin -permissions * -days 365

# 2. Copy the generated token
# 3. Configure server to use tokens.json file
```

**For Windows Users - Using Tokens with gfl Client:**
After creating a token, Windows users can set it as an environment variable for the gfl client:

```powershell
# Set the token in your PowerShell session
$env:GOFLUX_TOKEN_LITE = "your-generated-token-here"

# Or set it permanently for your user account
[System.Environment]::SetEnvironmentVariable("GOFLUX_TOKEN_LITE", "your-generated-token-here", "User")
```

**Note for End Users:** If you're not the administrator, you may need to ask your system administrator to:
1. Create a token for you using `gfl-admin create`
2. Provide you with the token value
3. Assign appropriate permissions for your use case

Then you can set the `GOFLUX_TOKEN_LITE` environment variable with the token they provide.

### Adding Users
```bash
# Create specific users with limited permissions
.\gfl-admin.exe create -user backup_service -permissions upload -days 30
.\gfl-admin.exe create -user web_app -permissions upload,download -days 90
.\gfl-admin.exe create -user monitoring -permissions list -days 7
```

### Token Maintenance
```bash
# Review all tokens
.\gfl-admin.exe list

# Remove compromised or unused tokens
.\gfl-admin.exe revoke old_token_id

# Create replacement tokens before expiration
.\gfl-admin.exe create -user alice -permissions * -days 90
```

## Integration with Server

Configure the server to use authentication by setting the `tokens_file` in your config:

**goflux.json:**
```json
{
  "server": {
    "tokens_file": "tokens.json"
  }
}
```

When tokens are configured, the server will require authentication for all operations and display a green confirmation message on startup.

## Troubleshooting

**Token file not found:**
- Create your first token with `gfl-admin create`
- Check file path and permissions

**Permission denied:**
- Verify token has required permissions
- Check token hasn't expired or been revoked
- Use `gfl-admin list` to verify token status

**Server not enforcing authentication:**
- Ensure `tokens_file` is set in server config
- Verify tokens.json exists and is readable
- Check server startup messages for authentication status
