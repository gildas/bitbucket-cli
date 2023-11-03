/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bitbucket.org/gildas_cherruel/bb/cmd"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
