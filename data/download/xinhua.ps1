# Download xinhua dictionary dataset
# Required for: dictctl fill-xinhua (fills character.json with Chinese meanings)
# Source: https://github.com/pwxcoo/chinese-xinhua
# License: Open data
# Size: ~26 MB

param(
    [string]$TargetDir = "$PSScriptRoot\..\xinhua"
)

# Mirror: raw GitHub content for word.json
$XinhuaUrl = "https://raw.githubusercontent.com/pwxcoo/chinese-xinhua/master/data/word.json"
$TargetFile = "$TargetDir\word.json"

New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null

Write-Host "Downloading xinhua word.json (~26MB)..."
Invoke-WebRequest -Uri $XinhuaUrl -OutFile $TargetFile

Write-Host "Done. Xinhua data at: $TargetFile"
Write-Host "Now run: go run ./cmd/dictctl fill-xinhua resources/character.json $TargetFile"
