$installDir = "$HOME\.sw\bin"
if (!(Test-Path $installDir)) { 
    New-Item -ItemType Directory -Force -Path $installDir 
}

$url = "https://github.com/zenith-sw/aws-role-switcher/releases/latest/download/sw_windows_amd64.exe"
$dest = Join-Path $installDir "sw.exe"

Write-Host "`nDownloading AWS STS Role Switcher for Windows..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $url -OutFile $dest

$path = [Environment]::GetEnvironmentVariable("Path", "User")
if ($path -notlike "*$installDir*") {
    $newPath = "$path;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path += ";$installDir"
    Write-Host ""
    Write-Host "Added $installDir to your User Path."
}

Write-Host "Starting to initiate configuration..."
& "$dest" init

Write-Host "`n-----" -ForegroundColor Gray
Write-Host "Installation Complete!" -ForegroundColor Green
Write-Host "Register your first role using 'sw add'"
Write-Host "`n-----" -ForegroundColor Gray