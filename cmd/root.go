package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/branch"
	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootOptions describes the options for the application
type RootOptions struct {
	ConfigFile     string `mapstructure:"-"`
	LogDestination string `mapstructure:"-"`
	ProfileName    string `mapstructure:"-"`
	Verbose        bool   `mapstructure:"-"`
	Debug          bool   `mapstructure:"-"`
}

// CmdOptions contains the options for the application
var CmdOptions RootOptions

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     APP,
	Version: Version(),
	Short:   "BitBucket Command Line Interface",
	Long: `BitBucket Command Line Interface is a tool to manage your BitBucket.
You can manage your pull requests, issues, profiles, etc.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(context context.Context) error {
	return rootCmd.ExecuteContext(context)
}

func init() {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&CmdOptions.ConfigFile, "config", "", "config file (default is .env, "+filepath.Join(configDir, "bitbucket", "config-cli.yml"))
	rootCmd.PersistentFlags().StringVarP(&CmdOptions.ProfileName, "profile", "p", "", "Profile to use. Overrides the default profile")
	rootCmd.PersistentFlags().StringVarP(&CmdOptions.LogDestination, "log", "l", "", "Log destination (stdout, stderr, file, none), overrides LOG_DESTINATION environment variable")
	rootCmd.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "logs are written at DEBUG level, overrides DEBUG environment variable")
	rootCmd.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "Verbose mode, overrides VERBOSE environment variable")

	rootCmd.AddCommand(profile.Command)
	rootCmd.AddCommand(project.Command)
	rootCmd.AddCommand(branch.Command)
	rootCmd.AddCommand(commit.Command)
	rootCmd.AddCommand(pullrequest.Command)

	rootCmd.SilenceUsage = true // Do not show usage when an error occurs
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log := logger.Must(logger.FromContext(rootCmd.Context()))

	if len(CmdOptions.LogDestination) > 0 {
		log.ResetDestinations(CmdOptions.LogDestination)
	}
	if CmdOptions.Debug {
		log.SetFilterLevel(logger.DEBUG)
	}

	log.Infof(strings.Repeat("-", 80))
	log.Infof("Starting %s v%s (%s)", APP, Version(), runtime.GOARCH)
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
		if err := profile.Profiles.Load(); err != nil {
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
			profile.Current = profile.Profiles.Current()
		}
		log.Infof("Profile: %s", profile.Current)
	}
}
