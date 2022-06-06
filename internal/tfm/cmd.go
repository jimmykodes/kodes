package tfm

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/internal/command"
)

var tfmCmd = &cobra.Command{
	Use:   "tfm",
	Short: "run a line by line transformation on input text",
	Args:  cobra.MinimumNArgs(1),
	Run:   run,
}

var output string

func New() *command.Command {
	return command.New(tfmCmd, flagInit)
}

func flagInit() error {
	tfmCmd.Flags().StringVarP(&output, "output", "o", "", "output file (default stdout)")

	return nil
}

func run(cmd *cobra.Command, args []string) {
	var (
		transformation string
		input          io.Reader
		out            io.Writer
		err            error
	)
	switch len(args) {
	case 1:
		transformation = args[0]
		input = os.Stdin
	case 2:
		transformation = args[0]
		input, err = os.Open(args[1])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tfm error: %s", err)
		return
	}

	if output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(output)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tfm error: %s", err)
		return
	}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Fprintf(out, transformation+"\n", line); err != nil {
			fmt.Fprintf(os.Stderr, "tfm error: %s", err)
		}
	}
}
