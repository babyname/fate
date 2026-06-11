# Download chinese-poetry dataset
# Required for: server poetry index, dbinit import-poetry
# Source: https://github.com/chinese-poetry/chinese-poetry (MIT)
# Size: ~354 MB

param(
    [string]$TargetDir = "$PSScriptRoot\..\chinese-poetry"
)

$RepoUrl = "https://github.com/chinese-poetry/chinese-poetry.git"

if (Test-Path "$TargetDir\.git") {
    Write-Host "Updating existing chinese-poetry checkout..."
    Push-Location $TargetDir
    git pull --depth=1
    Pop-Location
} else {
    Write-Host "Cloning chinese-poetry (shallow, ~354MB)..."
    if (Test-Path $TargetDir) {
        Remove-Item -Recurse -Force $TargetDir
    }
    git clone --depth=1 --filter=blob:none $RepoUrl $TargetDir
}

Write-Host "Done. Dataset at: $TargetDir"
Write-Host "Now run: go run ./cmd/dbinit import-poetry   (imports into DB)"
Write-Host "  or:  go run ./cmd/dictctl import-poetry     (generates poem_entries.json)"
