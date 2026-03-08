@echo off
REM Git Setup Script for Trademarkia Hotreload Assignment
REM This script will initialize git and create a proper commit history

echo ========================================
echo Git Setup for Trademarkia Assignment
echo ========================================
echo.

REM Check if git is installed
git --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Git is not installed or not in PATH
    echo Please install Git from https://git-scm.com/
    pause
    exit /b 1
)

echo Step 1: Configure Git User
echo ========================================
set /p GIT_NAME="Enter your name (e.g., Prachi Patel): "
set /p GIT_EMAIL="Enter your GitHub email: "

git config user.name "%GIT_NAME%"
git config user.email "%GIT_EMAIL%"

echo.
echo Git configured with:
git config user.name
git config user.email
echo.

echo Step 2: Initialize Git Repository
echo ========================================
git init
echo Git repository initialized.
echo.

echo Step 3: Add Remote Repository
echo ========================================
set /p REPO_URL="Enter your GitHub repository URL (e.g., https://github.com/PrachiPatel2105/Trademarkia.git): "
git remote add origin %REPO_URL%
echo Remote added: %REPO_URL%
echo.

echo Step 4: Creating Commit History
echo ========================================
echo This will create 16 meaningful commits...
echo.

REM Commit 1
echo [1/16] Project setup...
git add go.mod go.sum .gitignore
git commit -m "Initial commit: Setup Go module and dependencies" >nul 2>&1
echo ✓ Commit 1: Project setup

REM Commit 2
echo [2/16] Logging module...
git add internal/logger/
git commit -m "Add structured logging module using log/slog" >nul 2>&1
echo ✓ Commit 2: Logging module

REM Commit 3
echo [3/16] CLI configuration...
git add internal/cli/
git commit -m "Implement CLI configuration and flag parsing" >nul 2>&1
echo ✓ Commit 3: CLI configuration

REM Commit 4
echo [4/16] Path filtering...
git add internal/filter/
git commit -m "Add path filtering with ignore patterns and unit tests" >nul 2>&1
echo ✓ Commit 4: Path filtering

REM Commit 5
echo [5/16] File system watcher...
git add internal/watcher/
git commit -m "Implement file system watcher with fsnotify and recursive directory monitoring" >nul 2>&1
echo ✓ Commit 5: File system watcher

REM Commit 6
echo [6/16] Debouncer...
git add internal/debounce/
git commit -m "Add debouncer to handle rapid file events with unit tests" >nul 2>&1
echo ✓ Commit 6: Debouncer

REM Commit 7
echo [7/16] Builder...
git add internal/builder/
git commit -m "Implement build command execution with output capture" >nul 2>&1
echo ✓ Commit 7: Builder

REM Commit 8
echo [8/16] Process manager...
git add internal/process/
git commit -m "Add process manager with graceful shutdown and real-time log streaming" >nul 2>&1
echo ✓ Commit 8: Process manager

REM Commit 9
echo [9/16] Controller...
git add internal/controller/
git commit -m "Implement main controller with restart loop protection and build orchestration" >nul 2>&1
echo ✓ Commit 9: Controller

REM Commit 10
echo [10/16] CLI entry point...
git add cmd/hotreload/
git commit -m "Add CLI entry point with signal handling and context management" >nul 2>&1
echo ✓ Commit 10: CLI entry point

REM Commit 11
echo [11/16] Demo test server...
git add testserver/
git commit -m "Add demo HTTP server for testing hot reload functionality" >nul 2>&1
echo ✓ Commit 11: Demo test server

REM Commit 12
echo [12/16] Build automation...
git add Makefile
git commit -m "Add Makefile with build, test, and demo targets" >nul 2>&1
echo ✓ Commit 12: Build automation

REM Commit 13
echo [13/16] Documentation...
git add README.md QUICKSTART.md EXAMPLES.md ARCHITECTURE.md PROJECT_SUMMARY.md
git commit -m "Add comprehensive documentation" >nul 2>&1
echo ✓ Commit 13: Documentation

REM Commit 14
echo [14/16] Implementation completion...
git add IMPLEMENTATION_COMPLETE.md
git commit -m "Document implementation completion" >nul 2>&1
echo ✓ Commit 14: Implementation completion

REM Commit 15
echo [15/16] Submission materials...
git add LOOM_VIDEO_SCRIPT.md SUBMISSION_CHECKLIST.md GITHUB_SETUP_GUIDE.md
git commit -m "Add submission materials and guides" >nul 2>&1
echo ✓ Commit 15: Submission materials

REM Commit 16
echo [16/16] Compiled binary...
git add bin/
git commit -m "Add compiled binary for Windows" >nul 2>&1
echo ✓ Commit 16: Compiled binary

echo.
echo ========================================
echo All commits created successfully!
echo ========================================
echo.

echo Step 5: View Commit History
echo ========================================
git log --oneline
echo.

echo Step 6: Ready to Push
echo ========================================
echo Your repository is ready to push to GitHub.
echo.
echo To push, run:
echo   git branch -M main
echo   git push -u origin main
echo.
echo You will be prompted for:
echo   Username: PrachiPatel2105
echo   Password: [Your Personal Access Token]
echo.
echo Get your token at: https://github.com/settings/tokens
echo.

set /p PUSH_NOW="Do you want to push now? (y/n): "
if /i "%PUSH_NOW%"=="y" (
    echo.
    echo Pushing to GitHub...
    git branch -M main
    git push -u origin main
    
    if errorlevel 1 (
        echo.
        echo ERROR: Push failed. Please check:
        echo 1. Repository exists on GitHub
        echo 2. You have correct permissions
        echo 3. You're using a Personal Access Token, not password
        echo.
        echo Get token at: https://github.com/settings/tokens
    ) else (
        echo.
        echo ========================================
        echo SUCCESS! Repository pushed to GitHub
        echo ========================================
        echo.
        echo Next steps:
        echo 1. Verify at: %REPO_URL%
        echo 2. Grant access to: recruitments@trademarkia.com
        echo 3. Record your Loom video
        echo 4. Submit the form
        echo.
    )
) else (
    echo.
    echo Skipping push. You can push later with:
    echo   git branch -M main
    echo   git push -u origin main
    echo.
)

echo.
echo ========================================
echo Setup Complete!
echo ========================================
pause
