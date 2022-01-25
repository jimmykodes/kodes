package cmd

import "github.com/jimmykodes/kodes/internal/ts"

func init() {
	tsCmd := ts.New()
	tsCmd.Init()
	rootCmd.AddCommand(tsCmd.Cmd())
}
