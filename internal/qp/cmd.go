package qp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/internal/command"
	"github.com/jimmykodes/kodes/internal/util"
)

var qpCmd = &cobra.Command{
	Use:   "qp",
	Short: "Parse query params into json",
	Long: `Parse query params into a json object.

- Query params can be provided from an input file using --input or from stdin if none provided.
- Results can be sent to a file using --output or to stdout if none provided.
- Query Params should not be prepended with the ?. Though this won't break anything, the first key will just
have the ? prepended to it.
- Keys that show up more than once will have their values represented in an array.
- Unless --raw is provided, values will be parsed into best-guess data types.

Example
input: num=1&bool=false&str=string&rep=1&rep=2
output: {"bool":false,"num":1,"rep":[1,2],"str":"string"}
`,
	Args: cobra.MaximumNArgs(3),
	RunE: runE,
}

func New() *command.Command {
	return command.New(qpCmd, flagInit)
}

var (
	input, output string
	raw           bool
)

func flagInit() error {
	qpCmd.Flags().StringVarP(&input, "input", "i", "", "Input file with query params (default: stdin)")
	qpCmd.Flags().StringVarP(&output, "output", "o", "", "Output file for json (default: stdout)")
	qpCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Do not convert data types")
	return nil
}

func runE(_ *cobra.Command, args []string) error {
	var (
		src  io.ReadCloser
		dest io.WriteCloser
		err  error
	)
	if input == "" {
		src = os.Stdin
	} else {
		src, err = os.Open(input)
		if err != nil {
			return err
		}
	}
	if output == "" {
		dest = os.Stdout
	} else {
		dest, err = os.Create(args[1])
		if err != nil {
			return err
		}
	}

	defer func() {
		src.Close()
		dest.Close()
	}()

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	data = bytes.TrimSpace(data)

	u := url.URL{RawQuery: string(data)}

	final := make(map[string]interface{})
	for k, values := range u.Query() {
		if n := len(values); n == 1 {
			final[k] = convert(values[0])
		} else {
			v := make([]interface{}, n)
			for i, value := range values {
				v[i] = convert(value)
			}
			final[k] = v
		}
	}
	return json.NewEncoder(dest).Encode(final)
}

func convert(v string) interface{} {
	if raw {
		return v
	}
	return util.Convert(v)
}
