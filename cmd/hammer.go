package cmd

import "github.com/jimmykodes/kodes/internal/hammer"

func init() {
	hammerCmd := hammer.New()
	hammerCmd.Init()
	rootCmd.AddCommand(hammerCmd.Cmd())
}
