# build.ps1
param(
    [string]$version = "v1.0.0"
)

Write-Host "Starting build version $version..."

$platforms = @(
    @{GOOS="windows"; GOARCH="amd64"; Ext=".exe"},
    @{GOOS="linux";   GOARCH="amd64"; Ext=""}
)

foreach ($p in $platforms) {
    $outDir = "bin\$($p.GOOS)"
    if (-not (Test-Path $outDir)) { New-Item -ItemType Directory -Path $outDir }

    $outFile = "$outDir\monitoring$p.Ext"
    Write-Host "Building $($p.GOOS)/$($p.GOARCH) -> $outFile"
    
    $env:GOOS = $p.GOOS
    $env:GOARCH = $p.GOARCH
    go build -o $outFile main.go
}

Write-Host "Build complete!"

# Optional: Commit & Tag
git add .
git commit -m "Build version $version"
git tag $version
git push origin main --tags

Write-Host "Version $version pushed to GitHub"
