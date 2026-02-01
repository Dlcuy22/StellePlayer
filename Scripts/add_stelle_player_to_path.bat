@echo off
setlocal enabledelayedexpansion

:: =========================================================================
:: Endmi PATH Configuration Script
:: This script adds the directory containing endmi.exe to your User PATH.
:: =========================================================================

set "DEFAULT_PATH=C:\Program Files\StellePlayer"
set "TARGET_DIR="

echo Checking for endmi.exe in default location: %DEFAULT_PATH%

if exist "%DEFAULT_PATH%\StellePlayer.exe" (
    set "TARGET_DIR=%DEFAULT_PATH%"
    echo Found endmi.exe in default location.
) else (
    echo.
    echo [!] Could not find endmi.exe in %DEFAULT_PATH%
    echo.
    set /p "INPUT_PATH=Please paste the absolute path to the DIRECTORY containing endmi.exe: "
    
    :: Remove quotes if any
    set "INPUT_PATH=!INPUT_PATH:"=!"
    
    if exist "!INPUT_PATH!\StellePlayer.exe" (
        set "TARGET_DIR=!INPUT_PATH!"
    ) else (
        echo.
        echo [ERROR] endmi.exe was not found in: !INPUT_PATH!
        echo Please make sure you provide the folder path, not the executable path.
        pause
        exit /b 1
    )
)

echo.
echo Adding "!TARGET_DIR!" to User PATH...

:: Use PowerShell to safely append to User PATH (prevents 1024 char truncation issue with setx)
powershell -Command ^
    "[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH', 'User') + ';' + '%TARGET_DIR%', 'User')"

if %ERRORLEVEL% equ 0 (
    echo.
    echo [SUCCESS] StellePlayer has been added to your PATH.
    echo [INFO] Please restart your terminal/IDE for the changes to take effect.
) else (
    echo.
    echo [ERROR] Failed to update PATH.
)

pause
exit /b 0