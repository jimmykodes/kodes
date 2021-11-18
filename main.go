package main

import (
	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/cmd"
)

func main() {
	cobra.CheckErr(cmd.Execute())
}
