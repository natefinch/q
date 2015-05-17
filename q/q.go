package q

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"npf.io/q/q/log"
)

// ConfigDir reports the directory where Q should store its data.  The default
// is $HOME/.config/Q/ on *nixes and %LOCALAPPDATA%\Q\ on Windows.  The default
// may be overridden using the Q_CONFIG_DIR environment variable.
var ConfigDir = getConfigDir()

var (
	configFile = filepath.Join(ConfigDir, "q-config.toml")
)

func getConfigDir() string {
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

type Config struct {
	PluginDir string
}

func ReadConfig() (Config, error) {
	// init with defaults
	c := defaultConfig()
	meta, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		return c, fmt.Errorf("failed to read config file %q: %v", configFile, err)
	}
	if len(meta.Undecoded()) > 0 {
		log.Verbose("Unexpected options in Q config file: %v", meta.Undecoded())
	}
	return c, nil
}

func defaultConfig() Config {
	return Config{
		PluginDir: filepath.Join(ConfigDir, "plugins"),
	}
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
