# Download Unihan database
# Required for: dictctl import-unihan
# Source: https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip
# License: Unicode License (free)
# Size: ~2 MB compressed, ~49 MB extracted

param(
    [string]$TargetDir = "$PSScriptRoot\..\raw\unihan"
)

$UnihanUrl = "https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip"
$ZipPath = "$env:TEMP\Unihan.zip"

Write-Host "Downloading Unihan.zip..."
Invoke-WebRequest -Uri $UnihanUrl -OutFile $ZipPath

Write-Host "Extracting to $TargetDir..."
if (Test-Path $TargetDir) {
    Remove-Item -Recurse -Force $TargetDir
}
New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
Expand-Archive -Path $ZipPath -DestinationPath $TargetDir -Force

Remove-Item $ZipPath
Write-Host "Done. Unihan data at: $TargetDir"
Write-Host "Now run: go run ./cmd/dictctl import-unihan $TargetDir\Unihan_IRGSources.txt"
