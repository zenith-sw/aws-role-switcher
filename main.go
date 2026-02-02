package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type RoleInfo struct {
	Arn    string `yaml:"role_arn"`
	Region string `yaml:"region,omitempty"`
}

type ProfileConfig struct {
	Profile     string              `yaml:"profile"`
	Name        string              `yaml:"name"`
	Duration    int32               `yaml:"duration"`
	AssumeRoles map[string]RoleInfo `yaml:"assume_roles"`
}

var (
	configDir  string
	configPath string
)

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".aws")
	configPath = filepath.Join(configDir, "config.yaml")
}

var rootCmd = &cobra.Command{
	Use:   "sw",
	Short: "AWS Role Switcher CLI",
}

// --- Helper Functions ---

func loadConfigAndCheckUser() ([]ProfileConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read settings, please run 'sw init' first")
	}

	var cfg []ProfileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %v", err)
	}

	if len(cfg) == 0 || strings.TrimSpace(cfg[0].Name) == "" {
		return nil, fmt.Errorf("User name is missing. Please run 'sw init' again.")
	}

	return cfg, nil
}

func getSTSCredentials(targetAlias string) (*sts.AssumeRoleOutput, string) {
    profiles, err := loadConfigAndCheckUser()
    if err != nil { log.Fatal(err) }

    role, ok := profiles[0].AssumeRoles[targetAlias]
    if !ok { log.Fatalf("'%s' role not found", targetAlias) }

    os.Unsetenv("AWS_ACCESS_KEY_ID")
    os.Unsetenv("AWS_SECRET_ACCESS_KEY")
    os.Unsetenv("AWS_SESSION_TOKEN")

    ctx := context.TODO()
    var opts []func(*config.LoadOptions) error
    if role.Region != "" {
        opts = append(opts, config.WithRegion(role.Region))
    }

    cfg, err := config.LoadDefaultConfig(ctx, opts...)
    if err != nil { log.Fatalf("Unable to load SDK config: %v", err) }
    
    stsClient := sts.NewFromConfig(cfg)
    sessionName := fmt.Sprintf("%s-%s", profiles[0].Name, targetAlias)
    
    res, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
        RoleArn:         &role.Arn,
        RoleSessionName: &sessionName,
        DurationSeconds: &profiles[0].Duration,
    })
    if err != nil { log.Fatalf("Role switching failed: %v", err) }

    return res, role.Region
}

// --- Commands ---

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("Configuration already exists:", configPath)
			fmt.Println("You can manually edit the configuration file")
			return
		}

		var userName string
		for {
			fmt.Print("Enter your name (for session identification): ")
			fmt.Scanln(&userName)
			userName = strings.TrimSpace(userName)
			if userName != "" {
				break
			}
			fmt.Println("Name is required. Please try again.")
		}

		os.MkdirAll(configDir, 0755)
		defaultConfig := []ProfileConfig{{
			Profile:     "default",
			Name:        userName,
			Duration:    3600,
			AssumeRoles: make(map[string]RoleInfo),
		}}

		data, _ := yaml.Marshal(&defaultConfig)
		os.WriteFile(configPath, data, 0644)
		fmt.Println("\nInitialization completed successfully.")
		fmt.Println("Next step: Register your first role using 'sw add'")
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update a role ARN",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := loadConfigAndCheckUser()
		if err != nil {
			fmt.Println(err)
			return
		}

		var profileAlias, roleArn, region string
		var isUpdated bool

		fmt.Println("Add new AWS Role Profile\n")
		
		for {
			fmt.Print("Enter profile alias (e.g., dev, prod): ")
			fmt.Scanln(&profileAlias)
			profileAlias = strings.TrimSpace(profileAlias)
			if profileAlias != "" {
				break
			}
		}
		
		if _, exists := profiles[0].AssumeRoles[profileAlias]; exists {
			fmt.Printf("Profile '%s' already exists. Overwrite? (y/n): ", profileAlias)
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Cancelled.")
				return
			}
			isUpdated = true
		}

		for {
			fmt.Print("Enter Role ARN: ")
			fmt.Scanln(&roleArn)
			roleArn = strings.TrimSpace(roleArn)
			if roleArn != "" {
				break
			}
		}

		fmt.Print("Enter Region (optional): ")
		fmt.Scanln(&region)
		region = strings.TrimSpace(region)

		if profiles[0].AssumeRoles == nil {
			profiles[0].AssumeRoles = make(map[string]RoleInfo)
		}
		profiles[0].AssumeRoles[profileAlias] = RoleInfo{Arn: roleArn, Region: region}

		updatedData, _ := yaml.Marshal(&profiles)
		os.WriteFile(configPath, updatedData, 0644)
		
		if isUpdated {
			fmt.Printf("Successfully updated profile: [%s]\n", profileAlias)
		} else {
			fmt.Printf("Successfully registered new profile: [%s]\n", profileAlias)
		}
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup [alias]",
	Short: "Copy credentials to clipboard",
	Long: "Retrieves temporary credentials and copies them to your clipboard for manual use in your current terminal session.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetAlias := args[0]
        res, region := getSTSCredentials(targetAlias)

		var output string
    	if runtime.GOOS == "windows" {
			output = fmt.Sprintf(
                "$env:AWS_ACCESS_KEY_ID=\"%s\"; $env:AWS_SECRET_ACCESS_KEY=\"%s\"; $env:AWS_SESSION_TOKEN=\"%s\"",
                *res.Credentials.AccessKeyId, *res.Credentials.SecretAccessKey, *res.Credentials.SessionToken,
            )
            if region != "" {
                output += fmt.Sprintf("; $env:AWS_REGION=\"%s\"", region)
            }
		} else {
			output = fmt.Sprintf(
				"export AWS_ACCESS_KEY_ID=%s\nexport AWS_SECRET_ACCESS_KEY=%s\nexport AWS_SESSION_TOKEN=%s",
				*res.Credentials.AccessKeyId, *res.Credentials.SecretAccessKey, *res.Credentials.SessionToken,
			)
			if region != "" {
				output += fmt.Sprintf("\nexport AWS_REGION=%s", region)
			}
		}

		clipboard.WriteAll(output)
    	fmt.Printf("[%s] credentials copied! Paste and execute in your terminal.\n", targetAlias)
	},
}

var runCmd = &cobra.Command{
    Use:   "run [alias] -- [command]",
    Short: "Run a command with temporary credentials",
    Long: `Executes a command in a secure subshell with credentials automatically injected, leaving your main environment clean.
		Usage Example:
		- sw run {profile} -- {your_command}
  		- sw run {profile} -- {your_script.sh}`,
    Args:  cobra.MinimumNArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        res, region := getSTSCredentials(args[0])
        
        commandName := args[1]
        commandArgs := args[2:]

        c := exec.Command(commandName, commandArgs...)
        
        env := os.Environ()
        env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *res.Credentials.AccessKeyId))
        env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *res.Credentials.SecretAccessKey))
        env = append(env, fmt.Sprintf("AWS_SESSION_TOKEN=%s", *res.Credentials.SessionToken))
        env = append(env, fmt.Sprintf("AWS_REGION=%s", region))
        
        c.Env = env
        c.Stdout = os.Stdout
        c.Stderr = os.Stderr
        c.Stdin = os.Stdin

        fmt.Printf("Running: %s %v\n", commandName, commandArgs)
        if err := c.Run(); err != nil {
            log.Fatalf("Execution failed: %v", err)
        }
    },
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered role profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := loadConfigAndCheckUser()
		if err != nil {
			fmt.Println(err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ALIAS\tREGION\tROLE ARN")
		fmt.Fprintln(w, "-----\t------\t--------")

		for alias, info := range profiles[0].AssumeRoles {
			region := info.Region
			if region == "" {
				region = "(default)"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", alias, region, info.Arn)
		}
		w.Flush()
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [profile]",
	Short: "Remove a role profile",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := loadConfigAndCheckUser()
		if err != nil {
			fmt.Println(err)
			return
		}

		var alias string
		if len(args) > 0 {
			alias = args[0]
		} else {
			fmt.Print("Enter profile alias to delete: ")
			fmt.Scanln(&alias)
		}

		if _, exists := profiles[0].AssumeRoles[alias]; !exists {
			fmt.Printf("Profile '%s' not found.\n", alias)
			return
		}

		delete(profiles[0].AssumeRoles, alias)
		updatedData, _ := yaml.Marshal(&profiles)
		os.WriteFile(configPath, updatedData, 0644)

		fmt.Printf("Successfully deleted profile: [%s]\n", alias)
	},
}

func main() {
	rootCmd.AddCommand(initCmd, setupCmd, runCmd, addCmd, listCmd, deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}