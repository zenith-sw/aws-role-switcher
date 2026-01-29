if (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Error "Please run it with administrator privileges."
    exit
}

$installDir = "$HOME\.sw\bin"
if (!(Test-Path $installDir)) { New-Item -ItemType Directory -Force -Path $installDir }
$dest = "$installDir\sw.exe"

$url = "https://github.com/zenith-sw/aws-role-switcher/releases/latest/download/sw_windows_amd64.exe"
$dest = "C:\Windows\System32\sw.exe"

Write-Host "Downloading AWS STS Role Switcher for Windows..."
Invoke-WebRequest -Uri $url -OutFile $dest

$path = [Environment]::GetEnvironmentVariable("Path", "User")
if ($path -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$path;$installDir", "User")
    $env:Path += ";$installDir"
    Write-Host "Added $installDir to your User Path."
}

Write-Host "AWS STS Role Switcher is completely installed!"
& sw init

Write-Host "\n-----"
Write-Host "Installation Complete!"
Write-Host "---"
Write-Host "Next Step:"
Write-Host "1. Open your config file: vi ~/.sw/config.yaml"
Write-Host "2. Register your IAM Role ARNs under 'assume_roles'."
Write-Host "3. Run 'sw setup {profile}' to get your credentials!"
Write-Host "---"
Write-Host "Example: sw setup dev"