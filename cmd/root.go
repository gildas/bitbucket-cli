package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/artifact"
	"bitbucket.org/gildas_cherruel/bb/cmd/branch"
	"bitbucket.org/gildas_cherruel/bb/cmd/cache"
	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/component"
	"bitbucket.org/gildas_cherruel/bb/cmd/gpg-key"
	"bitbucket.org/gildas_cherruel/bb/cmd/issue"
	"bitbucket.org/gildas_cherruel/bb/cmd/pipeline"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	sshkey "bitbucket.org/gildas_cherruel/bb/cmd/ssh-key"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// RootOptions describes the options for the application
type RootOptions struct {
	ConfigFile     string         `mapstructure:"-"`
	LogDestination string         `mapstructure:"-"`
	ProfileName    string         `mapstructure:"-"`
	OutputFormat   flags.EnumFlag `mapstructure:"-"`
	DryRun         bool           `mapstructure:"-"`
	Verbose        bool           `mapstructure:"-"`
	Debug          bool           `mapstructure:"-"`
	StopOnError    bool           `mapstructure:"-"`
	WarnOnError    bool           `mapstructure:"-"`
	IgnoreErrors   bool           `mapstructure:"-"`
}

// CmdOptions contains the options for the application
var CmdOptions RootOptions

// RootCmd represents the base command when called without any subcommands
var RootCmd = NewRootCommand()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(context context.Context) error {
	return RootCmd.ExecuteContext(context)
}

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
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
	initializeRootCommand(root)
	return root
}

func initializeRootCommand(root *cobra.Command) {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	// Global flags
	CmdOptions = newRootOptions()

	root.PersistentFlags().StringVar(&CmdOptions.ConfigFile, "config", core.GetEnvAsString("BB_CONFIG", ""), "config file (default is .env, "+filepath.Join(configDir, "bitbucket", "config-cli.yml"))
	root.PersistentFlags().StringVarP(&CmdOptions.ProfileName, "profile", "p", core.GetEnvAsString("BB_PROFILE", ""), "Profile to use. Overrides the default profile")
	root.PersistentFlags().StringVarP(&CmdOptions.LogDestination, "log", "l", "", "Log destination (stdout, stderr, file, none), overrides LOG_DESTINATION environment variable")
	root.PersistentFlags().BoolVar(&CmdOptions.DryRun, "dry-run", false, "Dry run, the command will not modify anything but tell what it would do. \nAlso known as --noop, --what-if, or --whatif")
	root.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "logs are written at DEBUG level, overrides DEBUG environment variable")
	root.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "Verbose mode, overrides VERBOSE environment variable")
	root.PersistentFlags().VarP(&CmdOptions.OutputFormat, "output", "o", "Output format (json, yaml, table). Overrides the default output format of the profile")
	root.PersistentFlags().BoolVar(&CmdOptions.StopOnError, "stop-on-error", false, "Stop on error")
	root.PersistentFlags().BoolVar(&CmdOptions.WarnOnError, "warn-on-error", false, "Warn on error")
	root.PersistentFlags().BoolVar(&CmdOptions.IgnoreErrors, "ignore-errors", false, "Ignore errors")
	root.MarkFlagsMutuallyExclusive("stop-on-error", "warn-on-error", "ignore-errors")
	_ = root.MarkFlagFilename("config")
	_ = root.MarkFlagFilename("log")
	_ = root.RegisterFlagCompletionFunc("profile", profile.ValidProfileNames)
	_ = root.RegisterFlagCompletionFunc(CmdOptions.OutputFormat.CompletionFunc("output"))
	root.PersistentFlags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		switch name {
		case "noop", "dryrun", "whatif", "what-if":
			name = "dry-run"
		}
		return pflag.NormalizedName(name)
	})

	root.AddCommand(artifact.Command)
	root.AddCommand(profile.Command)
	root.AddCommand(project.Command)
	root.AddCommand(branch.Command)
	root.AddCommand(commit.Command)
	root.AddCommand(component.Command)
	root.AddCommand(issue.Command)
	root.AddCommand(pipeline.Command)
	root.AddCommand(pullrequest.Command)
	root.AddCommand(repository.Command)
	root.AddCommand(user.Command)
	root.AddCommand(workspace.Command)
	root.AddCommand(gpgkey.Command)
	root.AddCommand(sshkey.Command)
	root.AddCommand(cache.Command)

	root.SilenceUsage = true // Do not show usage when an error occurs
}

func newRootOptions() RootOptions {
	return RootOptions{
		OutputFormat: flags.EnumFlag{
			Allowed: []string{"csv", "json", "yaml", "table", "tsv"},
			Value:   core.GetEnvAsString("BB_OUTPUT_FORMAT", ""),
		},
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log := logger.Must(logger.FromContext(RootCmd.Context()))

	if len(CmdOptions.LogDestination) > 0 {
		log.ResetDestinations(CmdOptions.LogDestination)
	}
	if CmdOptions.Debug {
		log.SetFilterLevel(logger.DEBUG)
	}

	log.Infof("%s", strings.Repeat("-", 80))
	log.Infof("Starting %s v%s (%s)", RootCmd.Name(), RootCmd.Version, runtime.GOARCH)
	log.Infof("Log Destination: %s", log)

	viper.SetConfigType("yaml")
	if len(CmdOptions.ConfigFile) > 0 { // Use config file from the flag.
		viper.SetConfigFile(CmdOptions.ConfigFile)
	} else if configDir, _ := os.UserConfigDir(); len(configDir) > 0 {
		viper.AddConfigPath(filepath.Join(configDir, "bitbucket"))
		viper.SetConfigName("config-cli.yml")
	} else { // Old fashion configuration file in $HOME
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".bitbucket-cli")
	}

	// Read the config file
	err := viper.ReadInConfig()
	if verr, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warnf("Config file not found: %s", verr)
		if len(CmdOptions.ProfileName) > 0 {
			log.Fatalf("Profile %s not found (missing config file)", CmdOptions.ProfileName)
			fmt.Fprintf(os.Stderr, "Profile %s not found (missing config file)\n", CmdOptions.ProfileName)
			os.Exit(1)
		}
	} else if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
		fmt.Fprintf(os.Stderr, "Failed to read config file: %s\n", err)
		os.Exit(1)
	} else {
		log.Infof("Config File: %s", viper.ConfigFileUsed())
		if err := profile.Profiles.Load(RootCmd.Context()); err != nil {
			log.Fatalf("Failed to load profiles: %s", err)
			fmt.Fprintf(os.Stderr, "Failed to load profiles: %s\n", err)
			os.Exit(1)
		}
		if len(CmdOptions.ProfileName) > 0 {
			var found bool

			if profile.Current, found = profile.Profiles.Find(CmdOptions.ProfileName); !found {
				log.Fatalf("Profile %s not found", CmdOptions.ProfileName)
				fmt.Fprintf(os.Stderr, "Profile %s not found in %s\n", CmdOptions.ProfileName, viper.ConfigFileUsed())
				os.Exit(1)
			}
		} else {
			profile.Current = profile.Profiles.Current(RootCmd.Context())
		}
		log.Record("profile", profile.Current).Infof("Current Profile: %s", profile.Current)
	}
}
