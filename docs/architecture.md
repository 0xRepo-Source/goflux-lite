# GoFlux Lite Architecture

This document provides an overview of how the three GoFlux Lite components interact with each other.

## Component Interaction Diagram

```mermaid
flowchart TD
    Admin[gfl-admin<br/>Token Management]
    Server[gfl-server<br/>File Server]
    Client[gfl<br/>File Operations]
    TokenFile[(tokens.json)]
    ConfigFile[(goflux.json)]
    Storage[(data/ directory)]
    MetaDir[(.goflux-meta/)]
    EnvVar[GOFLUX_TOKEN_LITE]
    
    Admin -->|creates/lists/revokes| TokenFile
    Server -->|reads| TokenFile
    Server -->|reads| ConfigFile
    Server -->|stores files| Storage
    Server -->|manages sessions| MetaDir
    Client -->|reads| ConfigFile
    Client -->|HTTP API| Server
    Client -.->|optional| EnvVar
```

## Authentication Flow

```mermaid
sequenceDiagram
    participant A as gfl-admin
    participant TF as tokens.json
    participant S as gfl-server
    participant C as gfl
    
    Note over A,C: 1. Setup Phase
    A->>TF: create token
    TF-->>A: token created
    
    Note over A,C: 2. Server Startup
    S->>TF: read tokens
    TF-->>S: token data loaded
    S->>S: enable authentication
    
    Note over A,C: 3. Client Operations
    C->>C: read GOFLUX_TOKEN_LITE or config
    C->>S: API request + token
    S->>S: validate token
    S-->>C: authenticated response
## File Upload Flow

```mermaid
sequenceDiagram
    participant C as gfl client
    participant S as gfl-server
    participant FS as File System
    participant SS as Session Store
    
    Note over C,SS: Chunked Upload Process
    
    C->>S: POST /upload chunk 1
    S->>SS: create/update session
    SS-->>S: session updated
    S->>FS: store chunk
    FS-->>S: chunk stored
    S-->>C: chunk accepted
    
    C->>S: POST /upload chunk 2
    S->>SS: update session
    SS-->>S: session updated
    S->>FS: store chunk
    FS-->>S: chunk stored
    S-->>C: chunk accepted
    
    Note over C,SS: more chunks...
    
    C->>S: POST /upload final chunk
    S->>SS: complete session
    SS-->>S: session completed
    S->>FS: assemble final file
    FS-->>S: file complete
    S->>SS: cleanup session
    SS-->>S: session deleted
    S-->>C: upload complete
```

## Resume Workflow

```mermaid
sequenceDiagram
    participant C as gfl client
    participant S as gfl-server
    participant SS as Session Store
    
    Note over C,SS: Upload Interrupted Scenario
    
    C->>S: POST /upload chunk 1
    S->>SS: create session
    S-->>C: success
    
    C->>S: POST /upload chunk 2
    S->>SS: update session
    S-->>C: success
    
    Note over C: Network interruption!
    
    C->>S: GET /upload/status?path=file.txt
    S->>SS: get session
    SS-->>S: session data
    S-->>C: missing chunks [3,4,5,6]
    
    C->>S: POST /upload chunk 3
    S->>SS: update session
    S-->>C: success
    
    Note over C,SS: Resume continues...
```

## Component Responsibilities

### gfl-admin Token Management

```mermaid
flowchart TD
    Admin[gfl-admin] --> Create[Create Tokens]
    Admin --> List[List Tokens]
    Admin --> Revoke[Revoke Tokens]
    
    Create --> Validate[Validate Permissions]
    Create --> Generate[Generate Secure Token]
    Create --> Store[Store in tokens.json]
    
    List --> Read[Read tokens.json]
    List --> Format[Format Output]
    
    Revoke --> Find[Find Token]
    Revoke --> Mark[Mark as Revoked]
    Revoke --> Save[Save to tokens.json]
```

### gfl-server File Server

```mermaid
flowchart TD
    Server[gfl-server] --> Auth[Authentication]
    Server --> Upload[Upload Handler]
    Server --> Download[Download Handler]
    Server --> ListOp[List Handler]
    
    Auth --> LoadTokens[Load tokens.json]
    Auth --> ValidateReq[Validate Requests]
    Auth --> Challenge[Challenge-Response]
    
    Upload --> Chunks[Manage Chunks]
    Upload --> Sessions[Track Sessions]
    Upload --> Assemble[Assemble Files]
    
    Download --> ReadFile[Read File]
    Download --> Stream[Stream Response]
    
    ListOp --> ScanDir[Scan Directory]
    ListOp --> FilterPaths[Filter Paths]
```

### gfl Client

```mermaid
flowchart TD
    Client[gfl] --> Put[put command]
    Client --> Get[get command]
    Client --> Ls[ls command]
    
    Put --> ChunkFile[Split into Chunks]
    Put --> Upload[Upload Chunks]
    Put --> CheckStatus[Check Status]
    Put --> Resume[Resume if Needed]
    
    Get --> Request[Request File]
    Get --> Receive[Receive Data]
    Get --> SaveLocal[Save Locally]
    
    Ls --> ListReq[List Request]
    Ls --> ParseResp[Parse Response]
    Ls --> Display[Display Results]
```

## Security Architecture

```mermaid
flowchart TD
    subgraph Auth[Authentication Layer]
        Token[Token-based Auth]
        Challenge[Challenge-Response]
        Permissions[Permission System]
    end
    
    subgraph Network[Network Security]
        HTTPS[Optional HTTPS/TLS]
        Headers[Authorization Headers]
    end
    
    subgraph FileSystem[File System Security]
        PathSanitize[Path Sanitization]
        Traversal[Traversal Protection]
        Isolation[Storage Isolation]
    end
    
    ClientReq[Client Request] --> HTTPS
    HTTPS --> Token
    Token --> Permissions
    Permissions --> PathSanitize
    PathSanitize --> Isolation
    
    Challenge --> Headers
    Headers --> Traversal
```

## Data Flow Summary

1. **Setup**: Admin creates tokens using `gfl-admin`
2. **Server Start**: Server loads tokens and starts listening
3. **Client Auth**: Client uses token (env var or config) for authentication
4. **File Operations**: Client performs authenticated file operations
5. **Resume Support**: Interrupted uploads automatically resume
6. **Session Persistence**: Upload state survives server restarts

## Key Design Principles

- **Separation of Concerns**: Each tool has a single responsibility
- **Stateless Operations**: Server is stateless except for upload sessions
- **Security by Default**: Authentication required, path traversal protection
- **Resume Capability**: Robust handling of network interruptions
- **Simple Configuration**: Minimal setup required
- **Cross-Platform**: Works consistently across operating systems