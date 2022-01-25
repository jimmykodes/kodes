package ts

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/internal/command"
)

var tsCmd = &cobra.Command{
	Use:   "ts timestamps...",
	Short: "Convert an integer timestamp to an iso time string",
	Args:  cobra.MinimumNArgs(1),
	Run:   run,
}

func New() *command.Command {
	return command.New(tsCmd, flagInit)
}

func flagInit() error {
	// Add flags here
	return nil
}

func run(_ *cobra.Command, args []string) {
	for _, arg := range args {
		timestamp, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			fmt.Println("error processing timestamp", arg, "-", err)
			continue
		}
		switch len(arg) {
		case 10:
			t := time.Unix(timestamp, 0)
			fmt.Println(arg, "=>", t.Format(time.RFC3339))
		default:
			fmt.Println("invalid timestamp length")
		}
	}
}
