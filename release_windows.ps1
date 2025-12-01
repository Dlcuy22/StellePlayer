#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

$AppName = "StellePlayer"
$Version = "1.0.0"
$BinaryName = "Player.exe"
$BuildDir = "build"
$WindowsDir = "$BuildDir/windows"
$NsiScript = "scripts/build.nsi"

Write-Host "$AppName Build Script" -ForegroundColor Cyan


if (Test-Path $BuildDir) {
    Write-Host "Cleaning previous build artifacts..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force $BuildDir
}

# Create directories
New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
New-Item -ItemType Directory -Path $WindowsDir -Force | Out-Null

# Build Go binary
Write-Host "Building Go binary..." -ForegroundColor Green
go build -o "build/$BinaryName"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Go build failed!" -ForegroundColor Red
    exit 1
}

# Build NSIS installer
Write-Host "Building NSIS installer..." -ForegroundColor Green
$InstallerName = "${AppName}_Setup_${Version}.exe"

# Run makensis
& makensis $NsiScript
if ($LASTEXITCODE -ne 0) {
    Write-Host "NSIS build failed!" -ForegroundColor Red
    exit 1
}

# Move installer to build/windows
Move-Item -Force "scripts/$InstallerName" $WindowsDir

# Clean up binary
Remove-Item -Force "build/$BinaryName"

Write-Host ""
Write-Host "Build completed successfully!" -ForegroundColor Green