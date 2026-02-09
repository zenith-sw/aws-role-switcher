# AWS Role Switcher
**sw (AWS STS Role Switcher)** is the simplest way to manage and switch between AWS IAM Roles from your terminal.  

No more manually editing credentials or messy export commands.  
Just define your roles and switch in a second.  

**Instead of:**  
`aws sts assume-role --role-arn arn:aws:iam::... --role-session-name ...`
(and then manually exporting 3 variables)

**Just do:**  
`sw setup dev`

<br/>
<br/>

## Prerequisites
- **AWS CLI**: This tool **assumes that the AWS CLI is already configured** on your system.  
  If you haven't set it up yet, please run `aws configure` first.  

<br/>
<br/>

## Installation
**mac/linux:**  
```bash
curl -fsSL https://raw.githubusercontent.com/zenith-sw/aws-role-switcher/main/install.sh | bash
```

**Windows:**
```powershell
powershell -Command "irm https://raw.githubusercontent.com/zenith-sw/aws-role-switcher/main/install_windows.ps1 | iex"
```

<br/>
<br/>

## Quick Start
**Register Your Roles**  
Register Your Roles using `ars add` or you acn fill it menually.

Open the config file and add your IAM Role ARNs:
- macOS: vi ~/.aws/config.yaml
- Windows: notepad $HOME\.aws\config.yaml

**Assume Role**   
Now you can switch roles instantly:
```bash
sw setup {profile_alias}
```

<br/>
<br/>

## Uninstall
If you wish to remove `sw` from your system:

**macOS:**  
```bash
sudo rm /usr/local/bin/sw && rm -rf ~/.sw
```

**Windows:**
```powershell
Remove-Item "$HOME\.sw" -Recurse -Force
```

<br/>
<br/>

## Build from Sorce
**macOS:**  
```bash
GOOS=darwin GOARCH=arm64 go build -o sw_darwin_arm64 main.go
```

**Windows:**
```powershell
GOOS=windows GOARCH=amd64 go build -o sw_windows_amd64.exe main.go
```

**linuxOS:**  
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o sw_linux_amd64 main.go
```

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-s -w' -o sw_linux_arm64 main.go
```