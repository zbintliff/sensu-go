package check

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/sensu/sensu-go/cli"
	"github.com/sensu/sensu-go/cli/commands/hooks"
	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

// DeleteCommand adds a command that allows user to delete checks
func DeleteCommand(cli *cli.SensuCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete [NAME]",
		Short:        "delete checks given name",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no name is present print out usage
			if len(args) != 1 {
				cmd.Help()
				return nil
			}

			name := args[0]

			if skipConfirm, _ := cmd.Flags().GetBool("skip-confirm"); !skipConfirm {
				if ok, err := ConfirmDelete(name, cmd.OutOrStdout()); err != nil {
					return err
				} else if !ok {
					fmt.Fprintln(cmd.OutOrStdout(), "Canceled")
					return nil
				}
			}

			check := &types.CheckConfig{Name: name}
			err := cli.Client.DeleteCheck(check)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "OK")
			return nil
		},
		Annotations: map[string]string{
			// We want to be able to run this command regardless of whether the CLI
			// has been configured.
			hooks.ConfigurationRequirement: hooks.ConfigurationNotRequired,
		},
	}

	cmd.Flags().Bool("skip-confirm", false, "skip interactive confirmation prompt")

	return cmd
}

func ConfirmDelete(name string, stdout io.Writer) (bool, error) {
	confirmation := strings.ToUpper(name)

	// TODO: Colourize to emphaize destructive action
	message := `
Are you sure you would like to delete resource '` + name + `'?
Type '` + confirmation + `' to confirm.
	`

	stdout.Write([]byte(message))

	rl, err := readline.New("> ")
	if err != nil {
		return false, err
	}
	defer rl.Close()

	line, _ := rl.Readline()
	return confirmation == line, nil
}
