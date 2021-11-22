package strman

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jimmykodes/strman"
	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/internal/command"
)

var strmanCmd = &cobra.Command{
	Use:   "strman {camel | pascal | snake | kebab | screaming-snake | screaming-kebab} [inputs...]",
	Short: "Reformat input strings to the chosen output format",
	Long: `Reformat the input strings to the chosen output type.

Valid output types:
- camel
- pascal
- snake
- kebab
- screaming-snake
- screaming-kebab

If [inputs...] are empty, this will read from stdin.
Input format from stdin should be one input per line.

If inputs contain spaces, these will be parsed as separators and be replaced, ie:

kodes strman snake "this is a test" // returns this_is_a_test
`,
	Args: cobra.MinimumNArgs(1),
	Run:  run,
}

func New() *command.Command {
	return command.New(strmanCmd)
}

func run(_ *cobra.Command, args []string) {
	var inputs []string
	if len(args) == 1 {
		// no [inputs...] read from stdin and split on " "
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("error reading input", err)
			return
		}
		inputs = strings.Split(string(b), "\n")
	} else {
		inputs = args[1:]
	}
	for i, input := range inputs {
		// if the input contains a space, cast it to kebab case so strman can correctly parse it
		// since it is not expecting spaces in inputs.
		inputs[i] = strings.ReplaceAll(input, " ", "-")
	}
	var formatter func(string) string
	switch args[0] {
	case "camel":
		formatter = strman.ToCamel
	case "pascal":
		formatter = strman.ToPascal
	case "snake":
		formatter = strman.ToSnake
	case "kebab":
		formatter = strman.ToKebab
	case "screaming-snake":
		formatter = strman.ToScreamingSnake
	case "screaming-kebab":
		formatter = strman.ToScreamingKebab

	}
	for _, input := range inputs {
		fmt.Println(formatter(input))
	}
}
