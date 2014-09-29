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
	maine()
}

// maine is a separate function so it can be called from tests.
func maine() {
	q := &cobra.Command{
		Use:   "q is a command",
		Short: "q is a do-everything CLI tool for busy developers",
		Long: `Q is an omnipotent being in the Star Trek universe.  
q is a do-everything CLI tool for busy developers.`,
	}

	add := &cobra.Command{
		Use:   "add <text>",
		Short: "add a new task with the given text.",
		Long:  "Add a new task with the given text",
		Run:   add,
	}
	q.AddCommand(add)

	del := &cobra.Command{
		Use:   "del <id>",
		Short: "delete a task with the given id.",
		Long:  "Delete a task with the given id",
		Run:   del,
	}
	q.AddCommand(del)

	list := &cobra.Command{
		Use:   "list",
		Short: "list all existing tasks",
		Long:  "list all existing tasks",
		Run:   list,
	}
	q.AddCommand(list)

	if err := q.Execute(); err != nil {
		os.Exit(1)
	}
}

func add(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
	t := q.Todo{Title: args[0]}
	if err := q.Add(t); err != nil {
		fmt.Fprintf(os.Stderr, "Failed adding task: %s", err)
	}
}

func del(cmd *cobra.Command, args []string) {
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

func list(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
		return
	}
	if err := q.List(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed listing tasks: %s", err)
	}
}
