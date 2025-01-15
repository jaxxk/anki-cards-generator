@echo off
SET PROGRAM_NAME=poggers.exe
SET INSTALL_DIR=%USERPROFILE%\bin

REM Step 1: Build the program
echo Building the Go program...
go build -o %PROGRAM_NAME%
IF %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to build the program. Ensure Go is installed and properly set up.
    exit /b 1
)
echo Build successful.

REM Step 2: Create install directory if it doesn't exist
if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)

REM Step 3: Move the executable to the install directory (overwrite if exists)
echo Installing the program to %INSTALL_DIR%...
move /Y %PROGRAM_NAME% "%INSTALL_DIR%"
IF %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to move the program to %INSTALL_DIR%.
    exit /b 1
)
echo Installation complete: %PROGRAM_NAME% successfully updated in %INSTALL_DIR%.

REM Step 4: Add install directory to PATH (if not already added)
echo Ensuring %INSTALL_DIR% is in PATH...
echo %PATH% | find /I "%INSTALL_DIR%" >nul
IF ERRORLEVEL 1 (
    setx PATH "%PATH%;%INSTALL_DIR%"
    echo Added %INSTALL_DIR% to PATH. Restart your terminal to apply changes.
) ELSE (
    echo %INSTALL_DIR% is already in PATH.
)

REM Step 5: Set environment variables for logging if not already set
echo Checking environment variables for logging...
for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v LOG_LEVEL 2^>nul') do set LOG_LEVEL=%%B
for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v LOG_MODE 2^>nul') do set LOG_MODE=%%B

IF "%LOG_LEVEL%" == "INFO" IF "%LOG_MODE%" == "production" (
    echo Logging environment variables are already set.
) ELSE (
    echo Setting environment variables for production logging...
    setx LOG_LEVEL "INFO" /M
    setx LOG_MODE "production" /M
    IF %ERRORLEVEL% NEQ 0 (
        echo Error: Failed to set environment variables. Please check your permissions.
        exit /b 1
    )
    echo Environment variables set successfully.
)

echo All steps complete! You can now run %PROGRAM_NAME% from anywhere.
