## AWS Role Switcher
**sw (AWS STS Role Switcher)** is the simplest way to manage and switch between AWS IAM Roles from your terminal.

No more manually editing credentials or messy export commands.
Just define your roles and switch in a second.

**Instead of:**
`aws sts assume-role --role-arn arn:aws:iam::... --role-session-name ...` (and then manually exporting 3 variables)

**Just do:**
`sw setup dev`


---

### Installation
**macOS:**  
```bash
curl -fsSL https://raw.githubusercontent.com/zenith-sw/aws-role-switcher/main/install_mac.sh | bash
```

**Windows:**
```powershell
powershell -Command "irm https://raw.githubusercontent.com/zenith-sw/aws-role-switcher/main/install_windows.ps1 | iex"
```
<br/>

### Quick Start
**Configure**   
Open the config file and add your IAM Role ARNs:
- macOS: vi ~/.aws/config.yaml
- Windows: notepad $HOME\.aws\config.yaml

**Register Your Roles**  
Register Your Roles Fill in your ~/.aws/config.yaml as follows:

```YAML
assume_roles:
  dev:
    role_arn: "{your_role_arn}"
  prod:
    role_arn: "{your_role_arn}"
```

**Assume Role**   
Now you can switch roles instantly:
```bash
sw setup dev
```

<br/>

### Uninstall
If you wish to remove `sw` from your system:

**macOS:**  
```bash
sudo rm /usr/local/bin/sw && rm -rf ~/.sw
```

**Windows:**
```powershell
Remove-Item "$HOME\.sw" -Recurse -Force
```