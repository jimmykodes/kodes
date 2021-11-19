package cmd

import "github.com/jimmykodes/kodes/internal/qp"

func init() {
	qpCmd := qp.New()
	qpCmd.Init()
	rootCmd.AddCommand(qpCmd.Cmd())
}
