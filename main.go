package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	//"github.com/spf13/viper"

	"npf.io/q/q"
)

func main() {
	cmd := &cobra.Command{
		Short: "q is a do-everything CLI tool for busy developers",
		Long: `Q is an omnipotent being in the Star Trek universe.  
q is a do-everything CLI tool for busy developers.`,
	}

	addCmd(cmd)
	delCmd(cmd)
	listCmd(cmd)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func delCmd(base *cobra.Command) {
	del := &cobra.Command{
		Use:   "del <id>",
		Short: "delete a task",
		Long:  "Delete a task with the given id.",
	}
	base.AddCommand(del)

	del.Run = func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.Usage()
			return
		}

		if err := q.Delete(id); err != nil {
			fmt.Fprintf(os.Stderr, "Failed deleting task: %s", err)
		}
	}
}

func addCmd(base *cobra.Command) {
	add := &cobra.Command{
		Use:   "add <title>",
		Short: "create a new task",
		Long:  "Create a new task with the given title.",
	}
	base.AddCommand(add)
	desc := add.Flags().StringP("desc", "d", "", "Long-form description of the task")
	add.Run = func(cmd *cobra.Command, args []string) {
		// title is required
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		t := q.Todo{Title: args[0]}
		if desc != nil {
			t.Desc = *desc
		}
		if err := q.Add(t); err != nil {
			fmt.Fprintf(os.Stderr, "Failed adding task: %s", err)
		}
	}
}

func listCmd(base *cobra.Command) {
	list := &cobra.Command{
		Use:   "list",
		Short: "list all existing tasks",
		Long:  "list all existing tasks",
	}
	base.AddCommand(list)
	list.Run = func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Usage()
			return
		}
		if err := q.List(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed listing tasks: %s", err)
		}
	}
}
