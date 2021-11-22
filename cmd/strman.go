package cmd

import "github.com/jimmykodes/kodes/internal/strman"

func init() {
	strmanCmd := strman.New()
	strmanCmd.Init()
	rootCmd.AddCommand(strmanCmd.Cmd())
}
