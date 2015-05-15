package q

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigDir reports the directory where Q should store its data.  The default
// is $HOME/.config/Q/ on *nixes and %LOCALAPPDATA%\Q\ on Windows.  The default
// may be overridden using the Q_CONFIG_DIR environment variable.
func ConfigDir() string {
	if dir := os.Getenv("Q_CONFIG_DIR"); dir != "" {
		return dir
	}

	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "Q")
	} else {
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			filepath.Join(xdg, "Q")
		}
		return filepath.Join(os.Getenv("HOME"), ".config", "Q")
	}
}

func configFile() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

func ReadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(configFile())
	v.SetEnvPrefix("Q")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

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
			Short: cmd.Name + " <context>",
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
	return nil
}
