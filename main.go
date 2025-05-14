package main

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"toomani/cmd"
)

func main() {
	viper.AutomaticEnv()

	err := cmd.NewRootCmd().Execute()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
