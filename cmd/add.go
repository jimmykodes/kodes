package cmd

import "github.com/jimmykodes/kodes/internal/add"

func init() {
	addCmd := add.New()
	addCmd.Init()
	rootCmd.AddCommand(addCmd.Cmd())
}
