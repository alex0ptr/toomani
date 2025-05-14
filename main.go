package main

import (
	"github.com/alex0ptr/toomani/cmd"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	err := cmd.NewRootCmd().Execute()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
