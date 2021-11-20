package ctj

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/internal/command"
	"github.com/jimmykodes/kodes/internal/util"
)

var ctjCmd = &cobra.Command{
	Use:   "ctj",
	Short: "A brief description of the command",
	Run:   run,
}

func New() *command.Command {
	return command.New(ctjCmd, flagInit)
}

var (
	delimiter, input, output string
)

func flagInit() error {
	// Add flags here
	ctjCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Source file delimiter")
	ctjCmd.Flags().StringVarP(&input, "input", "i", "stdin", "Input source")
	ctjCmd.Flags().StringVarP(&output, "output", "o", "stdout", "Output source")
	return nil
}

func run(_ *cobra.Command, _ []string) {
	var (
		delim rune
		src   io.ReadCloser
		dest  io.WriteCloser
		err   error
	)
	if delimiter == "tab" {
		delim = '\t'
	} else {
		delim = rune(delimiter[0])
	}
	if input == "stdin" {
		src = os.Stdin
	} else {
		src, err = os.Open(input)
		if err != nil {
			fmt.Println("Error opening file", err)
			return
		}
	}
	if output == "stdout" {
		dest = os.Stdout
	} else {
		dest, err = os.Create(output)
		if err != nil {
			fmt.Println("Error creating output file", err)
			return
		}
	}
	defer func() {
		src.Close()
		dest.Close()
	}()

	reader := csv.NewReader(src)
	reader.Comma = delim
	reader.LazyQuotes = true
	encoder := json.NewEncoder(dest)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	obj := make(map[string]interface{})
	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Println("Error reading row", err)
			return
		}
		for i, h := range header {
			obj[h] = util.Convert(row[i])
		}
		if err := encoder.Encode(obj); err != nil {
			fmt.Println("Error writing JSON", err)
			return
		}
	}
}
