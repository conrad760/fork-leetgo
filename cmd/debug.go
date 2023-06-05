package cmd

import (
	"bytes"
	"os"

	"github.com/goccy/go-json"
	"github.com/spf13/cobra"

	"github.com/j178/leetgo/config"
	"github.com/j178/leetgo/leetcode"
)

var inspectCmd = &cobra.Command{
	Use:    "inspect",
	Short:  "Inspect LeetCode API, developer only",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := leetcode.NewClient()
		resp, err := c.Inspect(args[0])
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
		cmd.Print(buf.String())
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:    "whoami",
	Short:  "Print the current user",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cred, err := leetcode.ReadCredentials()
		if err != nil {
			return err
		}
		c := leetcode.NewClient(leetcode.WithCredentials(cred))
		user, err := c.GetUserStatus()
		if !user.IsSignedIn {
			return leetcode.ErrUserNotSignedIn
		}
		if err != nil {
			return err
		}
		cmd.Println(user.Whoami(c))
		return nil
	},
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Show debug info",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cwd, _ := os.Getwd()
		projectConfig, err := os.ReadFile(cfg.ProjectConfigFile())
		if err != nil {
			projectConfig = []byte("No project config file found")
		}
		cmd.Println("Leetgo version info  :")
		cmd.Println("```")
		cmd.Println(buildVersion())
		cmd.Println("```")
		cmd.Println("Global config dir    :", cfg.ConfigDir())
		cmd.Println("Global config file   :", cfg.GlobalConfigFile())
		cmd.Println("Project root         :", cfg.ProjectRoot())
		cmd.Println("Working dir          :", cwd)
		cmd.Println("Project config file  :", cfg.ProjectConfigFile())
		cmd.Println("Project configuration:")
		cmd.Println("```yaml")
		cmd.Println(string(projectConfig))
		cmd.Println("```")
		cmd.Println("Full configuration   :")
		cmd.Println("```yaml")
		_ = cfg.Write(cmd.OutOrStdout(), false)
		cmd.Println("```")
	},
}
