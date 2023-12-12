/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
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
