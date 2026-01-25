package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/vcs/git"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	gitVerbose bool
	gitDryRun  bool
)

// gitCmd represents the git command group
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git helper commands",
	Long: `Git helper commands for repository management.

Note: Most git operations are available as shell aliases.
These commands provide additional utilities.

Examples:
  acorn git info              # Show repo info
  acorn git contributors      # Show contributors
  acorn git find "bug fix"    # Find commits
  acorn git clean-branches    # Clean merged branches`,
}

// gitInfoCmd shows repo info
var gitInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show git repository info",
	Long: `Display current git repository information.

Shows branch, remote, and status counts.

Examples:
  acorn git info
  acorn git info -o json`,
	RunE: runGitInfo,
}

// gitContributorsCmd shows contributors
var gitContributorsCmd = &cobra.Command{
	Use:   "contributors",
	Short: "Show repository contributors",
	Long: `Display list of contributors with commit counts.

Examples:
  acorn git contributors
  acorn git contributors -o json`,
	Aliases: []string{"contrib"},
	RunE:    runGitContributors,
}

// gitFindCmd finds commits
var gitFindCmd = &cobra.Command{
	Use:   "find <search>",
	Short: "Find commits by message",
	Long: `Search for commits containing the given text in the message.

Examples:
  acorn git find "bug fix"
  acorn git find "refactor"`,
	Args: cobra.ExactArgs(1),
	RunE: runGitFind,
}

// gitCleanBranchesCmd cleans merged branches
var gitCleanBranchesCmd = &cobra.Command{
	Use:   "clean-branches",
	Short: "Clean merged branches",
	Long: `Remove local branches that have been merged.

Skips main, master, and develop branches.

Examples:
  acorn git clean-branches
  acorn git clean-branches --dry-run`,
	Aliases: []string{"cleanup"},
	RunE:    runGitCleanBranches,
}

// gitReposDirCmd shows repos directory
var gitReposDirCmd = &cobra.Command{
	Use:   "repos-dir",
	Short: "Show default repos directory",
	Long: `Display the default directory for git repositories.

Examples:
  acorn git repos-dir`,
	RunE: runGitReposDir,
}

func init() {
	vcsCmd.AddCommand(gitCmd)

	// Add subcommands
	gitCmd.AddCommand(gitInfoCmd)
	gitCmd.AddCommand(gitContributorsCmd)
	gitCmd.AddCommand(gitFindCmd)
	gitCmd.AddCommand(gitCleanBranchesCmd)
	gitCmd.AddCommand(gitReposDirCmd)

	// Persistent flags
	gitCmd.PersistentFlags().BoolVarP(&gitVerbose, "verbose", "v", false,
		"Show verbose output")
	gitCmd.PersistentFlags().BoolVar(&gitDryRun, "dry-run", false,
		"Show what would be done without executing")
}

func runGitInfo(cmd *cobra.Command, args []string) error {
	helper := git.NewHelper(gitVerbose)
	info := helper.GetInfo()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	// Table format
	if !info.IsRepo {
		fmt.Fprintf(os.Stdout, "%s Not a git repository\n", output.Warning("⚠"))
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Git Repository Info"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Fprintf(os.Stdout, "Branch: %s\n", output.Success(info.Branch))

	if info.Remote != "" {
		fmt.Fprintf(os.Stdout, "Remote: %s\n", info.Remote)
		fmt.Fprintf(os.Stdout, "URL:    %s\n", info.RemoteURL)
	}

	if info.Ahead > 0 || info.Behind > 0 {
		fmt.Fprintf(os.Stdout, "Status: ↑%d ↓%d\n", info.Ahead, info.Behind)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Staged:    %d\n", info.Staged)
	fmt.Fprintf(os.Stdout, "Modified:  %d\n", info.Modified)
	fmt.Fprintf(os.Stdout, "Untracked: %d\n", info.Untracked)

	return nil
}

func runGitContributors(cmd *cobra.Command, args []string) error {
	helper := git.NewHelper(gitVerbose)
	contributors, err := helper.GetContributors()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]interface{}{"contributors": contributors})
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Repository Contributors"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, c := range contributors {
		fmt.Fprintf(os.Stdout, "%6d  %s\n", c.Commits, c.Name)
	}

	return nil
}

func runGitFind(cmd *cobra.Command, args []string) error {
	helper := git.NewHelper(gitVerbose)
	commits, err := helper.FindCommits(args[0])
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string][]string{"commits": commits})
	}

	// Table format
	if len(commits) == 0 {
		fmt.Fprintf(os.Stdout, "No commits found matching: %s\n", args[0])
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Commits matching: "+args[0]))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, c := range commits {
		fmt.Fprintln(os.Stdout, c)
	}

	fmt.Fprintf(os.Stdout, "\nFound: %d commits\n", len(commits))
	return nil
}

func runGitCleanBranches(cmd *cobra.Command, args []string) error {
	helper := git.NewHelper(gitVerbose)
	deleted, err := helper.CleanMergedBranches(gitDryRun)
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		fmt.Fprintln(os.Stdout, "No merged branches to clean")
		return nil
	}

	if gitDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Branches that would be deleted:"))
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Deleted branches:"))
	}

	for _, b := range deleted {
		fmt.Fprintf(os.Stdout, "  %s %s\n", output.Success("✓"), b)
	}

	fmt.Fprintf(os.Stdout, "\nTotal: %d branches\n", len(deleted))
	return nil
}

func runGitReposDir(cmd *cobra.Command, args []string) error {
	helper := git.NewHelper(gitVerbose)
	fmt.Fprintln(os.Stdout, helper.GetReposDir())
	return nil
}
