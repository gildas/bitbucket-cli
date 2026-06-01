package main

import (
	"context"
	"os"

	"github.com/gildas/bitbucket-cli/cmd"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // .env is optional; error is expected when the file is absent
	if len(os.Getenv("LOG_DESTINATION")) == 0 {
		os.Setenv("LOG_DESTINATION", "nil")
	}
	log := logger.Create(APP)
	defer log.Flush()
	cmd.RootCmd.Use = APP
	cmd.RootCmd.Version = Version()
	err := cmd.Execute(log.ToContext(context.Background()))
	if err != nil {
		log.Fatalf("Failed to execute command", err)
		os.Exit(1)
	}
}
