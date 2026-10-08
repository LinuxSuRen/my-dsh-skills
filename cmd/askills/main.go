// Command askills embeds the skills of linuxsuren/my-dsh-skills and
// installs them for AI coding tools (dsh, codex, opencode and any
// agents-compatible tool).
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/linuxsuren/my-dsh-skills"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "askills",
		Short: "Install embedded agent skills for dsh, codex, opencode and agents-compatible tools",
	}
	cmd.AddCommand(newListCmd(), newInstallCmd(), newUninstallCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List embedded skills",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			metas, err := skills.All()
			if err != nil {
				return err
			}
			for _, m := range metas {
				fmt.Printf("%-24s %s\n", m.Name, m.Description)
			}
			return nil
		},
	}
}

func newInstallCmd() *cobra.Command {
	var (
		tool      string
		dir       string
		skillList []string
	)
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install skills to AI coding tools (default: shared ~/.agents/skills)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			targets, err := resolveTargets(tool, dir)
			if err != nil {
				return err
			}
			for name, dir := range targets {
				installed, err := skills.Install(dir, skillList)
				if err != nil {
					return err
				}
				fmt.Printf("installed %d skill(s) to %s (%s): %s\n",
					len(installed), dir, name, strings.Join(installed, ", "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&tool, "tool", "agents",
		fmt.Sprintf("target tool: %s (default agents, shared by dsh/codex/opencode)", strings.Join([]string{skills.ToolAgents, skills.ToolDSH, skills.ToolOpenCode, skills.ToolCodex, skills.ToolAll}, ", ")))
	cmd.Flags().StringVar(&dir, "dir", "", "custom target directory (overrides --tool)")
	cmd.Flags().StringSliceVar(&skillList, "skill", nil, "skill names to install (default: all)")
	return cmd
}

func newUninstallCmd() *cobra.Command {
	var (
		tool      string
		dir       string
		skillList []string
	)
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove installed skills",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			targets, err := resolveTargets(tool, dir)
			if err != nil {
				return err
			}
			for name, dir := range targets {
				removed, err := skills.Uninstall(dir, skillList)
				if err != nil {
					return err
				}
				fmt.Printf("removed %d skill(s) from %s (%s): %s\n",
					len(removed), dir, name, strings.Join(removed, ", "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&tool, "tool", "agents", "target tool, see install")
	cmd.Flags().StringVar(&dir, "dir", "", "custom target directory (overrides --tool)")
	cmd.Flags().StringSliceVar(&skillList, "skill", nil, "skill names to remove (default: all)")
	return cmd
}

func resolveTargets(tool, dir string) (map[string]string, error) {
	if dir != "" {
		return map[string]string{"custom": dir}, nil
	}
	return skills.TargetDirs(tool)
}
