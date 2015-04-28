package q

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Short: "q is a do-everything CLI tool for busy people",
	Long: `Q is an omnipotent being in the Star Trek universe.  
q is a do-everything CLI tool for busy people.`,
}

type Command struct {
	Name      string
	Context   string
	PluralCtx string
	Exec      func(args []string) bool
	Validate  func(args []string) bool
	Usage     string
	Short     string
	Long      string
}

func Add(cmd Command) error {
	var command *cobra.Command
	for _, c := range baseCmd.Commands() {
		if c.Name() == cmd.Name {
			command = c
			break
		}
	}
	if command == nil {
		command = &cobra.Command{
			Short: cmd + " <context>",
			Long:  fmt.Sprintf("%s something. Type q help %s for more details.", cmd, cmd),
		}
		baseCmd.AddCommand(command)
	}
	for _, c := range command.Commands() {
		if c.Name() == cmd.Context {
			return fmt.Errorf("Context %q already exists for command %q", cmd.Context, cmd.Name)
		}
	}
	ctxCmd := &cobra.Command{
		Use:   cmd.Usage,
		Short: cmd.Short,
		Long:  cmd.Long,
		Run: func(cc *cobra.Command, args []string) {
			if !cmd.Validate(args) {
				cc.Usage()
				return
			}
			if !cmd.Exec(args) {
				os.Exit(1)
			}
		},
	}

	command.AddCommand(ctxCmd)
}
