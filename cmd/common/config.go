package common

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initialize configure the logger and load the Viper Configuration
func Initialize(ctx context.Context, cmd *cobra.Command) (err error) {
	initializeLogger(ctx, cmd)
	return initializeConfiguration(ctx, cmd)
}

// initializeLogger configures the logger based on the command line flags and environment variables
func initializeLogger(ctx context.Context, cmd *cobra.Command) {
	log := logger.Must(logger.FromContext(ctx))

	// Use persistent flags instead
	if cmd.Root().PersistentFlags().Changed("log") {
		log.ResetDestinations(cmd.Root().PersistentFlags().Lookup("log").Value.String())
	}
	if cmd.Root().PersistentFlags().Changed("debug") && cmd.Root().PersistentFlags().Lookup("debug").Value.String() == "true" {
		log.SetFilterLevel(logger.DEBUG)
	}

	log.Infof("%s", strings.Repeat("-", 80))
	log.Infof("Starting %s v%s (%s)", cmd.Root().Name(), cmd.Root().Version, runtime.GOARCH)
	log.Infof("Log Destination: %s", log)
}

// initializeConfiguration loads the configuration file and profiles
func initializeConfiguration(ctx context.Context, cmd *cobra.Command) (err error) {
	log := logger.Must(logger.FromContext(ctx))

	viper.SetConfigType("yaml")
	if cmd.Root().PersistentFlags().Changed("config") {
		viper.SetConfigFile(cmd.Root().PersistentFlags().Lookup("config").Value.String())
	} else if configDir, _ := os.UserConfigDir(); len(configDir) > 0 {
		viper.AddConfigPath(filepath.Join(configDir, "bitbucket"))
		viper.SetConfigName("config-cli.yml")
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".bitbucket-cli")
	}

	err = viper.ReadInConfig()
	if verr, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warnf("Config file not found: %s", verr)
	} else if err != nil {
		return errors.Join(errors.New("Failed to read config file"), err)
	} else {
		log.Infof("Config File: %s", viper.ConfigFileUsed())
	}
	return nil
}
