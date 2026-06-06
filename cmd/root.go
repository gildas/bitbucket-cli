package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gildas/bitbucket-cli/cmd/artifact"
	"github.com/gildas/bitbucket-cli/cmd/branch"
	"github.com/gildas/bitbucket-cli/cmd/cache"
	"github.com/gildas/bitbucket-cli/cmd/commit"
	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/component"
	"github.com/gildas/bitbucket-cli/cmd/gpg-key"
	"github.com/gildas/bitbucket-cli/cmd/issue"
	"github.com/gildas/bitbucket-cli/cmd/pipeline"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/project"
	"github.com/gildas/bitbucket-cli/cmd/pullrequest"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	sshkey "github.com/gildas/bitbucket-cli/cmd/ssh-key"
	"github.com/gildas/bitbucket-cli/cmd/tag"
	"github.com/gildas/bitbucket-cli/cmd/user"
	"github.com/gildas/bitbucket-cli/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/spf13/cobra"
)

// RootOptions describes the options for the application
type RootOptions struct {
	ConfigFile     string          `mapstructure:"-"`
	LogDestination string          `mapstructure:"-"`
	ProfileName    string          `mapstructure:"-"`
	Repository     string          `mapstructure:"-"`
	Workspace      *flags.EnumFlag `mapstructure:"-"`
	OutputFormat   flags.EnumFlag  `mapstructure:"-"`
	DryRun         bool            `mapstructure:"-"`
	Verbose        bool            `mapstructure:"-"`
	Debug          bool            `mapstructure:"-"`
	StopOnError    bool            `mapstructure:"-"`
	WarnOnError    bool            `mapstructure:"-"`
	IgnoreErrors   bool            `mapstructure:"-"`
}

// CmdOptions contains the options for the application
var CmdOptions RootOptions

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Short: "BitBucket Command Line Interface",
	Long: `BitBucket Command Line Interface is a tool to manage your BitBucket.
You can manage your pull requests, issues, profiles, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bb requires a command:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(context context.Context) error {
	return RootCmd.ExecuteContext(context)
}

func init() {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	// Global flags
	CmdOptions.Workspace = flags.NewEnumFlagWithFunc(RootCmd, "", workspace.GetWorkspaceAllowedSlugs)
	CmdOptions.OutputFormat = flags.EnumFlag{Allowed: []string{"csv", "json", "yaml", "table", "tsv"}, Value: core.GetEnvAsString("BB_OUTPUT_FORMAT", "")}
	RootCmd.PersistentFlags().StringVar(&CmdOptions.ConfigFile, "config", core.GetEnvAsString("BB_CONFIG", ""), "config file (default is .env, "+filepath.Join(configDir, "bitbucket", "config-cli.yml"))
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.ProfileName, "profile", "p", core.GetEnvAsString("BB_PROFILE", ""), "Profile to use. Overrides the default profile")
	RootCmd.PersistentFlags().Var(CmdOptions.Workspace, "workspace", "Workspace to use. Overrides the default workspace of the profile. \nBy default, the workspace is determined from the git or profile configuration")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.Repository, "repository", "", "Repository to use. Overrides the default repository of the profile. \nBy default, the repository is determined from the git configuration")
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.LogDestination, "log", "l", "", "Log destination (stdout, stderr, file, none), overrides LOG_DESTINATION environment variable")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.DryRun, "dry-run", false, "Dry run, the command will not modify anything but tell what it would do. \nAlso known as --noop or --whatif")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.DryRun, "noop", false, "Dry run, the command will not modify anything but tell what it would do. \nAlso known as --dry-run or --whatif")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.DryRun, "whatif", false, "Dry run, the command will not modify anything but tell what it would do. \nAlso known as --dry-run or --noop")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "logs are written at DEBUG level, overrides DEBUG environment variable")
	RootCmd.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "Verbose mode, overrides VERBOSE environment variable")
	RootCmd.PersistentFlags().VarP(&CmdOptions.OutputFormat, "output", "o", "Output format (json, yaml, table). Overrides the default output format of the profile")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.StopOnError, "stop-on-error", false, "Stop on error")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	RootCmd.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
	_ = RootCmd.MarkFlagFilename("config")
	_ = RootCmd.MarkFlagFilename("log")
	_ = RootCmd.RegisterFlagCompletionFunc("profile", profile.ValidProfileNames)
	_ = RootCmd.RegisterFlagCompletionFunc(CmdOptions.OutputFormat.CompletionFunc("output"))
	_ = RootCmd.RegisterFlagCompletionFunc(CmdOptions.Workspace.CompletionFunc("workspace"))

	RootCmd.AddCommand(artifact.Command)
	RootCmd.AddCommand(profile.Command)
	RootCmd.AddCommand(project.Command)
	RootCmd.AddCommand(branch.Command)
	RootCmd.AddCommand(commit.Command)
	RootCmd.AddCommand(tag.Command)
	RootCmd.AddCommand(component.Command)
	RootCmd.AddCommand(issue.Command)
	RootCmd.AddCommand(pipeline.Command)
	RootCmd.AddCommand(pullrequest.Command)
	RootCmd.AddCommand(repository.Command)
	RootCmd.AddCommand(user.Command)
	RootCmd.AddCommand(workspace.Command)
	RootCmd.AddCommand(gpgkey.Command)
	RootCmd.AddCommand(sshkey.Command)
	RootCmd.AddCommand(cache.Command)

	RootCmd.SilenceUsage = true // Do not show usage when an error occurs
	cobra.OnInitialize(func() {
		if err := common.Initialize(RootCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize: %s\n", err)
			os.Exit(1)
		}
	})
}
