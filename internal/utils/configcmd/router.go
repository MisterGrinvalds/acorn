// Package configcmd provides a universal config subcommand router for any component
// with a .sapling/config/<component>/config.yaml.
//
// Components with a `files:` block get full path/source/generate/show functionality.
// Components without a `files:` block still get source info (config.yaml location).
//
// Usage:
//
//	opencodeCmd.AddCommand(configcmd.NewConfigRouter("opencode"))
package configcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/mistergrinvalds/acorn/internal/utils/configfile"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewConfigRouter returns a `config` cobra.Command with universal subcommands
// for managing a component's acorn-managed configuration.
func NewConfigRouter(component string) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  fmt.Sprintf("Commands for managing %s configuration files.", component),
	}

	pathCmd := &cobra.Command{
		Use:   "path",
		Short: "Show configuration file paths",
		Long: fmt.Sprintf(`Display target and generated paths for %s config files.

Examples:
  acorn ... %s config path
  acorn ... %s config path -o json`, component, component, component),
		Aliases: []string{"paths"},
		RunE:    runConfigPath(component),
	}

	sourceCmd := &cobra.Command{
		Use:   "source",
		Short: "Show configuration source info",
		Long: fmt.Sprintf(`Show where %s configuration is defined and generated.

Displays the component config.yaml path, the generated directory,
and symlink status for each managed file.

Examples:
  acorn ... %s config source
  acorn ... %s config source -o json`, component, component, component),
		RunE: runConfigSource(component),
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate configuration files",
		Long: fmt.Sprintf(`Re-render %s config files from the files: block in config.yaml.

Supports --dry-run to preview generated content without writing.

Examples:
  acorn ... %s config generate
  acorn ... %s config generate --dry-run`, component, component, component),
		Aliases: []string{"gen", "render"},
		RunE:    runConfigGenerate(component),
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Display generated configuration content",
		Long: fmt.Sprintf(`Show the current generated config file contents for %s.

Examples:
  acorn ... %s config show
  acorn ... %s config show -o json`, component, component, component),
		Aliases: []string{"cat"},
		RunE:    runConfigShow(component),
	}

	configCmd.AddCommand(pathCmd)
	configCmd.AddCommand(sourceCmd)
	configCmd.AddCommand(generateCmd)
	configCmd.AddCommand(showCmd)

	return configCmd
}

// componentConfigPath returns the path to a component's config.yaml.
func componentConfigPath(component string) string {
	root, _ := config.SaplingRoot()
	return filepath.Join(root, "config", component, "config.yaml")
}

// loadComponentFiles loads and parses the files: block from a component's config.yaml.
// Returns nil files (not an error) if the component has no files: block.
func loadComponentFiles(component string) ([]config.FileConfig, error) {
	configData, err := config.GetComponentConfig(component)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for %s: %w", component, err)
	}

	var base config.BaseConfig
	if err := yaml.Unmarshal(configData, &base); err != nil {
		return nil, fmt.Errorf("failed to parse config for %s: %w", component, err)
	}

	return base.Files, nil
}

// -- path subcommand --

// PathInfo holds path info for structured output.
type PathInfo struct {
	Target    string `json:"target" yaml:"target"`
	Expanded  string `json:"expanded" yaml:"expanded"`
	Generated string `json:"generated" yaml:"generated"`
	Format    string `json:"format" yaml:"format"`
}

func runConfigPath(component string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ioHelper := ioutils.IO(cmd)
		files, err := loadComponentFiles(component)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			if ioHelper.IsStructured() {
				return ioHelper.WriteOutput([]PathInfo{})
			}
			fmt.Fprintf(os.Stdout, "%s has no managed config files (no files: block in config.yaml)\n",
				output.Info(component))
			return nil
		}

		genDir, _ := config.GeneratedDir()

		paths := make([]PathInfo, 0, len(files))
		for _, f := range files {
			expanded := configfile.ExpandPath(f.Target)
			generated := ""
			if genDir != "" {
				generated = filepath.Join(genDir, component, filepath.Base(expanded))
			}
			paths = append(paths, PathInfo{
				Target:    f.Target,
				Expanded:  expanded,
				Generated: generated,
				Format:    f.Format,
			})
		}

		if ioHelper.IsStructured() {
			return ioHelper.WriteOutput(paths)
		}

		fmt.Fprintf(os.Stdout, "%s\n\n", output.Info(component+" Configuration Paths"))
		for _, p := range paths {
			fmt.Fprintf(os.Stdout, "  Target:     %s\n", p.Target)
			fmt.Fprintf(os.Stdout, "  Expanded:   %s\n", p.Expanded)
			if p.Generated != "" {
				fmt.Fprintf(os.Stdout, "  Generated:  %s\n", p.Generated)
			}
			fmt.Fprintf(os.Stdout, "  Format:     %s\n\n", p.Format)
		}

		return nil
	}
}

// -- source subcommand --

// SourceInfo holds source info for structured output.
type SourceInfo struct {
	Component    string     `json:"component" yaml:"component"`
	ConfigPath   string     `json:"config_path" yaml:"config_path"`
	HasFiles     bool       `json:"has_files" yaml:"has_files"`
	GeneratedDir string     `json:"generated_dir" yaml:"generated_dir"`
	Files        []FileLink `json:"files" yaml:"files"`
}

// FileLink holds per-file symlink status.
type FileLink struct {
	Target          string `json:"target" yaml:"target"`
	GeneratedPath   string `json:"generated_path" yaml:"generated_path"`
	SymlinkExists   bool   `json:"symlink_exists" yaml:"symlink_exists"`
	SymlinkValid    bool   `json:"symlink_valid" yaml:"symlink_valid"`
	GeneratedExists bool   `json:"generated_exists" yaml:"generated_exists"`
}

func runConfigSource(component string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ioHelper := ioutils.IO(cmd)

		configPath := componentConfigPath(component)
		genDir, _ := config.GeneratedDir()

		files, err := loadComponentFiles(component)
		if err != nil {
			// Config.yaml doesn't exist or can't be parsed — still show what we can
			files = nil
		}

		info := SourceInfo{
			Component:    component,
			ConfigPath:   configPath,
			HasFiles:     len(files) > 0,
			GeneratedDir: genDir,
			Files:        make([]FileLink, 0, len(files)),
		}

		for _, f := range files {
			expanded := configfile.ExpandPath(f.Target)
			genPath := ""
			if genDir != "" {
				genPath = filepath.Join(genDir, component, filepath.Base(expanded))
			}

			link := FileLink{
				Target:        expanded,
				GeneratedPath: genPath,
			}

			if genPath != "" {
				if _, statErr := os.Stat(genPath); statErr == nil {
					link.GeneratedExists = true
				}
			}

			if linkDest, linkErr := os.Readlink(expanded); linkErr == nil {
				link.SymlinkExists = true
				if genPath != "" {
					absLink, _ := filepath.Abs(linkDest)
					absGen, _ := filepath.Abs(genPath)
					link.SymlinkValid = absLink == absGen
				}
			}

			info.Files = append(info.Files, link)
		}

		if ioHelper.IsStructured() {
			return ioHelper.WriteOutput(info)
		}

		fmt.Fprintf(os.Stdout, "%s\n\n", output.Info(component+" Configuration Source"))
		fmt.Fprintf(os.Stdout, "  Component:      %s\n", info.Component)
		fmt.Fprintf(os.Stdout, "  Config YAML:    %s\n", info.ConfigPath)
		fmt.Fprintf(os.Stdout, "  Generated dir:  %s\n", info.GeneratedDir)

		if !info.HasFiles {
			fmt.Fprintf(os.Stdout, "\n  %s\n", output.Warning("No files: block defined"))
			return nil
		}

		fmt.Fprintln(os.Stdout)
		for _, fl := range info.Files {
			fmt.Fprintf(os.Stdout, "  Target: %s\n", fl.Target)
			if fl.GeneratedPath != "" {
				fmt.Fprintf(os.Stdout, "    Generated:  %s", fl.GeneratedPath)
				if fl.GeneratedExists {
					fmt.Fprintf(os.Stdout, " %s\n", output.Success("[exists]"))
				} else {
					fmt.Fprintf(os.Stdout, " %s\n", output.Warning("[missing]"))
				}
			}
			if fl.SymlinkExists {
				if fl.SymlinkValid {
					fmt.Fprintf(os.Stdout, "    Symlink:    %s\n", output.Success("valid"))
				} else {
					fmt.Fprintf(os.Stdout, "    Symlink:    %s\n", output.Warning("points elsewhere"))
				}
			} else {
				fmt.Fprintf(os.Stdout, "    Symlink:    %s\n", output.Warning("not linked"))
			}
			fmt.Fprintln(os.Stdout)
		}

		return nil
	}
}

// -- generate subcommand --

func runConfigGenerate(component string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ioHelper := ioutils.IO(cmd)
		files, err := loadComponentFiles(component)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			if ioHelper.IsStructured() {
				return ioHelper.WriteOutput([]*configfile.GeneratedFile{})
			}
			fmt.Fprintf(os.Stdout, "%s has no files: block — nothing to generate\n",
				output.Info(component))
			return nil
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		manager := configfile.NewManager(dryRun)

		results := make([]*configfile.GeneratedFile, 0, len(files))
		for _, f := range files {
			result, genErr := manager.GenerateFileForComponent(component, f)
			if genErr != nil {
				return fmt.Errorf("failed to generate %s: %w", f.Target, genErr)
			}
			results = append(results, result)
		}

		if ioHelper.IsStructured() {
			return ioHelper.WriteOutput(results)
		}

		if dryRun {
			fmt.Fprintf(os.Stdout, "%s\n\n", output.Info(component+" Config Generate (dry-run)"))
		} else {
			fmt.Fprintf(os.Stdout, "%s\n\n", output.Info(component+" Config Generate"))
		}

		for _, r := range results {
			status := output.Success("written")
			if dryRun {
				status = output.Warning("dry-run")
			}
			fmt.Fprintf(os.Stdout, "  %s %s\n", status, r.GeneratedPath)
			if dryRun {
				fmt.Fprintf(os.Stdout, "  Content:\n%s\n", r.Content)
			}
		}

		return nil
	}
}

// -- show subcommand --

func runConfigShow(component string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ioHelper := ioutils.IO(cmd)
		files, err := loadComponentFiles(component)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			if ioHelper.IsStructured() {
				return ioHelper.WriteOutput([]struct{}{})
			}
			fmt.Fprintf(os.Stdout, "%s has no files: block — nothing to show\n",
				output.Info(component))
			return nil
		}

		genDir, err := config.GeneratedDir()
		if err != nil {
			return fmt.Errorf("cannot determine generated directory: %w", err)
		}

		type ShowEntry struct {
			File    string `json:"file" yaml:"file"`
			Format  string `json:"format" yaml:"format"`
			Content string `json:"content" yaml:"content"`
		}

		entries := make([]ShowEntry, 0, len(files))
		for _, f := range files {
			expanded := configfile.ExpandPath(f.Target)
			genPath := filepath.Join(genDir, component, filepath.Base(expanded))

			content, readErr := os.ReadFile(genPath)
			if readErr != nil {
				entries = append(entries, ShowEntry{
					File:    genPath,
					Format:  f.Format,
					Content: fmt.Sprintf("(not found: %s)", readErr),
				})
				continue
			}
			entries = append(entries, ShowEntry{
				File:    genPath,
				Format:  f.Format,
				Content: string(content),
			})
		}

		if ioHelper.IsStructured() {
			return ioHelper.WriteOutput(entries)
		}

		for _, e := range entries {
			fmt.Fprintf(os.Stdout, "%s  [%s]\n", output.Info(e.File), e.Format)
			fmt.Fprintln(os.Stdout, e.Content)
		}

		return nil
	}
}
