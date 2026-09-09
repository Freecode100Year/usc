# USC (Universal Skill Compiler) One-Click Automated Installer for Windows
# Usage: powershell -ExecutionPolicy Bypass -File install.ps1
# Or:    irm https://raw.githubusercontent.com/Freecode100Year/usc/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "       USC (Universal Skill Compiler) One-Click Installer      " -ForegroundColor Cyan
Write-Host "================================================================" -ForegroundColor Cyan

# 1. Check Go environment
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Error "Go (golang) compiler is required to build USC from source. Please install Go from https://golang.org/dl/"
    exit 1
}

$repoRoot = $PSScriptRoot
if (-not $repoRoot -or -not (Test-Path "$repoRoot\go.mod")) {
    $repoRoot = (Get-Location).Path
}

Write-Host "[+] Building USC clean-room binary..." -ForegroundColor Yellow
Set-Location $repoRoot
go build -o bin\usc.exe .\cmd\usc
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to compile USC binary."
    exit 1
}

# 2. Determine global PATH target
$goBin = $env:GOPATH
if (-not $goBin) {
    $goBin = (go env GOPATH)
}
$targetDir = Join-Path $goBin "bin"
if (-not (Test-Path $targetDir)) {
    New-Item -ItemType Directory -Force -Path $targetDir | Out-Null
}

$targetExe = Join-Path $targetDir "usc.exe"
Copy-Item -Force "bin\usc.exe" $targetExe
Write-Host "[✓] Installed binary to: $targetExe" -ForegroundColor Green

# 3. Verify PATH availability
$pathEntries = [Environment]::GetEnvironmentVariable("PATH", "User") -split ';'
if ($pathEntries -notcontains $targetDir) {
    Write-Host "[!] Adding $targetDir to User PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("PATH", "$([Environment]::GetEnvironmentVariable('PATH', 'User'));$targetDir", "User")
    $env:PATH = "$env:PATH;$targetDir"
}

# 4. Auto-register USC skill into detected agents
$userHome = [Environment]::GetFolderPath("UserProfile")
$agySkillDir = Join-Path $userHome ".gemini\config\skills\usc"
if (Test-Path (Join-Path $userHome ".gemini")) {
    New-Item -ItemType Directory -Force -Path $agySkillDir | Out-Null
    $lines = @(
        "---",
        "name: usc",
        "description: Universal Skill Compiler (USC) - Compile, verify, and install skills from ClawHub, GitHub, or local source into local AI Agents.",
        "---",
        "",
        "# Universal Skill Compiler (USC)",
        "",
        "Use USC to compile and install any skill from ClawHub or GitHub into your local Agent:",
        "- 一键编译并安装 ClawHub 技能: usc install https://clawhub.ai/<author>/skills/<skill-name>",
        "- 一键编译并安装 GitHub 技能:  usc install https://github.com/<owner>/<repo>",
        "- 仅执行零信任安全编译与审计:   usc build <source-or-url>",
        "- 检查已支持的本地 Agent:        usc targets"
    )
    $lines | Set-Content -Path (Join-Path $agySkillDir "SKILL.md") -Encoding UTF8
    Write-Host "[✓] Registered USC capability skill into Antigravity CLI ($agySkillDir)" -ForegroundColor Green
}

Write-Host ""
Write-Host "================================================================" -ForegroundColor Green
Write-Host " [✓] USC Installation Complete!                                 " -ForegroundColor Green
Write-Host " You or your Agent can now run:                                 " -ForegroundColor Green
Write-Host "   usc install https://clawhub.ai/thesentitrader/skills/us-stocks-analysis" -ForegroundColor White
Write-Host "================================================================" -ForegroundColor Green

# 5. Show targets
& $targetExe targets
