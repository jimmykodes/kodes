package cmd

import "github.com/jimmykodes/kodes/internal/ctj"

func init() {
	ctjCmd := ctj.New()
	ctjCmd.Init()
	rootCmd.AddCommand(ctjCmd.Cmd())
}
