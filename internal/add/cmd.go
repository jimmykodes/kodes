package add

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jimmykodes/strman"
	"github.com/spf13/cobra"

	"github.com/jimmykodes/kodes/internal/command"
)

//go:embed templates
var fs embed.FS

var addCmd = &cobra.Command{
	Use:   "add command-name",
	Short: "Add a new command to the kodes base command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpls, err := template.ParseFS(fs, "templates/*.tmpl")
		if err != nil {
			return err
		}
		cmdName := args[0]

		packageName := strman.ToDelimited(cmdName, "")
		internalFilePath := filepath.Join("internal", packageName)

		if err := os.Mkdir(internalFilePath, 0755); err != nil {
			return fmt.Errorf("error creating internal package: %w", err)
		}

		short, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		tmplContext := &context{
			Package: packageName,
			Use:     strman.ToKebab(cmdName),
			Command: fmt.Sprintf("%sCmd", strman.ToCamel(cmdName)),
			Short:   short,
		}

		internalFile, err := os.Create(filepath.Join(internalFilePath, "cmd.go"))
		if err != nil {
			return err
		}
		defer internalFile.Close()
		if err := tmpls.ExecuteTemplate(internalFile, "internal.tmpl", tmplContext); err != nil {
			return err
		}

		cmdFile, err := os.Create(filepath.Join("cmd", fmt.Sprintf("%s.go", strman.ToSnake(cmdName))))
		if err != nil {
			return err
		}
		defer cmdFile.Close()

		if err := tmpls.ExecuteTemplate(cmdFile, "cmd.tmpl", tmplContext); err != nil {
			return err
		}

		return nil
	},
}

func New() *command.Command {
	return command.New(addCmd, flagInit)
}

func flagInit() error {
	addCmd.Flags().String("description", "A brief description of the command", "short command description")
	return nil
}

type context struct {
	Package string
	Use     string
	Command string
	Short   string
}
