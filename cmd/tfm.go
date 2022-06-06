package cmd

import "github.com/jimmykodes/kodes/internal/tfm"

func init() {
	tfmCmd := tfm.New()
	tfmCmd.Init()
	rootCmd.AddCommand(tfmCmd.Cmd())
}
