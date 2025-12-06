# GoFlux Lite Build Script for Windows

Write-Host "Building GoFlux Lite components..." -ForegroundColor Green

# Build server
Write-Host "Building server..." -ForegroundColor Yellow
go build -o bin/gfl-server.exe ./cmd/server
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed: server" -ForegroundColor Red
    exit 1
}

# Build client
Write-Host "Building client..." -ForegroundColor Yellow
go build -o bin/gfl.exe ./cmd/client
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed: client" -ForegroundColor Red
    exit 1
}

# Build admin
Write-Host "Building admin..." -ForegroundColor Yellow
go build -o bin/gfl-admin.exe ./cmd/admin
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed: admin" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Build successful!" -ForegroundColor Green
Write-Host ""
Write-Host "Binaries created:" -ForegroundColor Yellow
Write-Host "  gfl-server.exe - File server"
Write-Host "  gfl.exe - File client"
Write-Host "  gfl-admin.exe - Token management"
Write-Host ""
Write-Host "Quick start:" -ForegroundColor Yellow
Write-Host "  .\bin\gfl-server.exe -port 8080"
Write-Host "  .\bin\gfl-admin.exe create -user admin -permissions *"
Write-Host "  .\bin\gfl.exe update"
