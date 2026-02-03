package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/python"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	pythonDryRun  bool
	pythonVerbose bool
)

// pythonCmd represents the python command group
var pythonCmd = &cobra.Command{
	Use:   "python",
	Short: "Python development helpers",
	Long: `Helpers for Python development workflow with UV.

Provides commands for virtual environment management, dependency
management with UV, and development environment setup.

Examples:
  acorn python venv new myenv    # Create virtual environment
  acorn python init              # Initialize UV project
  acorn python add fastapi       # Add package
  acorn python sync              # Sync dependencies
  acorn python fastapi           # Setup FastAPI environment`,
	Aliases: []string{"py"},
}

// pythonVenvCmd is the parent for venv subcommands
var pythonVenvCmd = &cobra.Command{
	Use:   "venv",
	Short: "Virtual environment management",
	Long: `Manage Python virtual environments.

Examples:
  acorn python venv new .venv    # Create .venv
  acorn python venv list         # List available venvs`,
}

// pythonVenvNewCmd creates a new virtual environment
var pythonVenvNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a virtual environment",
	Long: `Create a new Python virtual environment using UV or python3 -m venv.

Uses UV if available (faster), otherwise falls back to python3.
Default name is ".venv".

Examples:
  acorn python venv new          # Create .venv
  acorn python venv new myenv    # Create myenv`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPythonVenvNew,
}

// pythonVenvListCmd lists virtual environments
var pythonVenvListCmd = &cobra.Command{
	Use:   "list",
	Short: "List virtual environments",
	Long: `List available virtual environments.

Checks both the default ENVS_LOCATION (~/.virtualenvs) and
the current directory for .venv.

Examples:
  acorn python venv list
  acorn python venv list -o json`,
	Aliases: []string{"ls"},
	RunE:    runPythonVenvList,
}

// pythonInitCmd initializes a new UV project
var pythonInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new Python project with UV",
	Long: `Initialize a new Python project using UV.

Creates pyproject.toml and basic project structure.

Examples:
  acorn python init              # Init in current directory
  acorn python init myproject    # Init new project directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPythonInit,
}

// pythonSyncCmd syncs dependencies
var pythonSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dependencies from pyproject.toml",
	Long: `Sync dependencies using UV.

Reads pyproject.toml and uv.lock to install dependencies.

Examples:
  acorn python sync`,
	RunE: runPythonSync,
}

// pythonAddCmd adds packages
var pythonAddCmd = &cobra.Command{
	Use:   "add <packages...>",
	Short: "Add packages to the project",
	Long: `Add packages using UV.

Examples:
  acorn python add fastapi
  acorn python add requests pandas numpy`,
	Args: cobra.MinimumNArgs(1),
	RunE: runPythonAdd,
}

// pythonRemoveCmd removes packages
var pythonRemoveCmd = &cobra.Command{
	Use:   "remove <packages...>",
	Short: "Remove packages from the project",
	Long: `Remove packages using UV.

Examples:
  acorn python remove requests
  acorn python remove pandas numpy`,
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	RunE:    runPythonRemove,
}

// pythonRunCmd runs a command in the project environment
var pythonRunCmd = &cobra.Command{
	Use:   "run <command> [args...]",
	Short: "Run a command in the project environment",
	Long: `Run a command using UV run.

Examples:
  acorn python run python main.py
  acorn python run pytest
  acorn python run ruff check .`,
	Args:               cobra.MinimumNArgs(1),
	RunE:               runPythonRun,
	DisableFlagParsing: true,
}

// pythonEnvCmd shows Python environment
var pythonEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show Python environment information",
	Long: `Display Python environment variables and configuration.

Examples:
  acorn python env
  acorn python env -o json`,
	RunE: runPythonEnv,
}

// pythonFastapiCmd sets up FastAPI environment
var pythonFastapiCmd = &cobra.Command{
	Use:   "fastapi [venv-name]",
	Short: "Setup FastAPI development environment",
	Long: `Create a virtual environment with FastAPI dependencies.

Installs:
  - fastapi, uvicorn, python-multipart
  - pytest, httpx, pytest-asyncio
  - ruff, python-dotenv

Examples:
  acorn python fastapi           # Create .venv with FastAPI
  acorn python fastapi myenv     # Create myenv with FastAPI`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPythonFastapi,
}

// pythonSetupCmd is the parent for setup subcommands
var pythonSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup development tools",
	Long: `Install common development tools.

Examples:
  acorn python setup ipython     # Install IPython with rich
  acorn python setup devtools    # Install ruff, mypy, pytest, etc.`,
}

// pythonSetupIPythonCmd installs IPython
var pythonSetupIPythonCmd = &cobra.Command{
	Use:   "ipython",
	Short: "Install IPython with rich output",
	Long: `Install IPython and rich for enhanced interactive Python.

Examples:
  acorn python setup ipython`,
	RunE: runPythonSetupIPython,
}

// pythonSetupDevtoolsCmd installs development tools
var pythonSetupDevtoolsCmd = &cobra.Command{
	Use:   "devtools",
	Short: "Install common development tools",
	Long: `Install development tools: ruff, mypy, pytest, pytest-cov, pre-commit.

Examples:
  acorn python setup devtools`,
	RunE: runPythonSetupDevtools,
}

func init() {

	// Venv subcommands
	pythonCmd.AddCommand(pythonVenvCmd)
	pythonVenvCmd.AddCommand(pythonVenvNewCmd)
	pythonVenvCmd.AddCommand(pythonVenvListCmd)

	// Main subcommands
	pythonCmd.AddCommand(pythonInitCmd)
	pythonCmd.AddCommand(pythonSyncCmd)
	pythonCmd.AddCommand(pythonAddCmd)
	pythonCmd.AddCommand(pythonRemoveCmd)
	pythonCmd.AddCommand(pythonRunCmd)
	pythonCmd.AddCommand(pythonEnvCmd)
	pythonCmd.AddCommand(pythonFastapiCmd)

	// Setup subcommands
	pythonCmd.AddCommand(pythonSetupCmd)
	pythonCmd.AddCommand(configcmd.NewConfigRouter("python"))
	pythonSetupCmd.AddCommand(pythonSetupIPythonCmd)
	pythonSetupCmd.AddCommand(pythonSetupDevtoolsCmd)

	// Persistent flags
	pythonCmd.PersistentFlags().BoolVar(&pythonDryRun, "dry-run", false,
		"Show what would be done without executing")
	pythonCmd.PersistentFlags().BoolVarP(&pythonVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runPythonVenvNew(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	name := ".venv"
	if len(args) > 0 {
		name = args[0]
	}

	info, err := helper.CreateVenv(name)
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	fmt.Fprintf(os.Stdout, "%s Virtual environment created!\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Name: %s\n", info.Name)
	fmt.Fprintf(os.Stdout, "  Path: %s\n", info.Path)
	if info.Python != "" {
		fmt.Fprintf(os.Stdout, "  Python: %s\n", info.Python)
	}
	fmt.Fprintf(os.Stdout, "  Created with: %s\n", info.CreatedBy)
	fmt.Fprintf(os.Stdout, "\nTo activate:\n")
	fmt.Fprintf(os.Stdout, "  source %s/bin/activate\n", name)

	return nil
}

func runPythonVenvList(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	venvs, err := helper.ListVenvs()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(venvs)
	}

	if len(venvs) == 0 {
		fmt.Fprintln(os.Stdout, "No virtual environments found.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Virtual Environments"))
	for _, v := range venvs {
		status := ""
		if v.Active {
			status = " " + output.Success("(active)")
		}
		fmt.Fprintf(os.Stdout, "  %s%s\n", v.Name, status)
		fmt.Fprintf(os.Stdout, "    Path: %s\n", v.Path)
		if v.Python != "" {
			fmt.Fprintf(os.Stdout, "    Python: %s\n", v.Python)
		}
	}

	return nil
}

func runPythonInit(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	if err := helper.InitProject(name); err != nil {
		return err
	}

	if name != "" {
		fmt.Fprintf(os.Stdout, "%s Python project initialized: %s\n", output.Success("✓"), name)
		fmt.Fprintf(os.Stdout, "\nNext steps:\n")
		fmt.Fprintf(os.Stdout, "  cd %s\n", name)
		fmt.Fprintf(os.Stdout, "  acorn python sync\n")
	} else {
		fmt.Fprintf(os.Stdout, "%s Python project initialized\n", output.Success("✓"))
		fmt.Fprintf(os.Stdout, "\nNext steps:\n")
		fmt.Fprintf(os.Stdout, "  acorn python sync\n")
	}

	return nil
}

func runPythonSync(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	fmt.Fprintln(os.Stdout, "Syncing dependencies...")
	if err := helper.Sync(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Dependencies synced\n", output.Success("✓"))
	return nil
}

func runPythonAdd(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	fmt.Fprintf(os.Stdout, "Adding packages: %v\n", args)
	if err := helper.Add(args...); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Packages added\n", output.Success("✓"))
	return nil
}

func runPythonRemove(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	fmt.Fprintf(os.Stdout, "Removing packages: %v\n", args)
	if err := helper.Remove(args...); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Packages removed\n", output.Success("✓"))
	return nil
}

func runPythonRun(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)
	return helper.Run(args...)
}

func runPythonEnv(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)
	info := helper.GetEnvInfo()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Python Environment"))
	if info.Version != "" {
		fmt.Fprintf(os.Stdout, "Python:        %s\n", info.Version)
	} else {
		fmt.Fprintf(os.Stdout, "Python:        %s\n", output.Warning("not found"))
	}
	if info.Pip != "" {
		fmt.Fprintf(os.Stdout, "Pip:           %s\n", info.Pip)
	}
	if info.UV != "" {
		fmt.Fprintf(os.Stdout, "UV:            %s\n", info.UV)
	} else {
		fmt.Fprintf(os.Stdout, "UV:            %s\n", output.Warning("not installed"))
	}
	if info.VirtualEnv != "" {
		fmt.Fprintf(os.Stdout, "Virtual Env:   %s\n", output.Success(info.VirtualEnv))
	} else {
		fmt.Fprintf(os.Stdout, "Virtual Env:   none active\n")
	}
	fmt.Fprintf(os.Stdout, "Envs Location: %s\n", info.EnvsLocation)

	return nil
}

func runPythonFastapi(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	name := ".venv"
	if len(args) > 0 {
		name = args[0]
	}

	info, err := helper.SetupFastAPI(name)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s FastAPI environment ready!\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Virtual env: %s\n", info.Path)
	fmt.Fprintf(os.Stdout, "\nTo activate:\n")
	fmt.Fprintf(os.Stdout, "  source %s/bin/activate\n", name)
	fmt.Fprintf(os.Stdout, "\nTo start development server:\n")
	fmt.Fprintf(os.Stdout, "  uvicorn main:app --reload\n")

	return nil
}

func runPythonSetupIPython(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	if err := helper.SetupIPython(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s IPython installed with rich output\n", output.Success("✓"))
	return nil
}

func runPythonSetupDevtools(cmd *cobra.Command, args []string) error {
	helper := python.NewHelper(pythonVerbose, pythonDryRun)

	if err := helper.SetupDevTools(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Development tools installed\n", output.Success("✓"))
	fmt.Fprintln(os.Stdout, "  Installed: ruff, mypy, pytest, pytest-cov, pre-commit")
	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "python",
		RegisterCmd: func() *cobra.Command { return pythonCmd },
	})
}
