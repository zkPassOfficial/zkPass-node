package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	optionNameDataDir        = "data-dir"
	optionNameSessionMax     = "session-max"
	optionNameSessionTimeout = "session-timeout"
	optionNameSessionLife    = "session-life"
)

type command struct {
	root    *cobra.Command
	config  *viper.Viper
	homeDir string
	cfgFile string
}

func WithHomeDir(dir string) func(c *command) {
	return func(c *command) {
		c.homeDir = dir
	}
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:           "zk-node",
			Short:         "zkPass Foundation zk-node",
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return c.initConfig()
			},
		},
	}
	for _, o := range opts {
		o(c)
	}

	if err := c.setHomeDir(); err != nil {
		return nil, err
	}

	c.initGlobalFlags()

	if err := c.initStartCmd(); err != nil {
		return nil, err
	}

	c.initVersionCmd()

	return c, nil
}

func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.root.Execute()
}

func init() {
	cobra.EnableCommandSorting = false
}

func (c *command) initConfig() (err error) {
	config := viper.New()
	configName := ".zkpass-node"
	if c.cfgFile != "" {
		config.SetConfigFile(c.cfgFile)
	} else {
		// Search config in home directory with name ".zkpass-node".
		config.AddConfigPath(c.homeDir)
		config.SetConfigName(configName)
	}

	// Environment
	config.AutomaticEnv() // read in environment variables that match
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if c.homeDir != "" && c.cfgFile == "" {
		c.cfgFile = filepath.Join(c.homeDir, configName+".yaml")
		fmt.Println("c.cfgFile:", c.cfgFile)
	}

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			return err
		}
	}
	c.config = config
	fmt.Println(config.AllSettings())
	return nil
}

func (c *command) initGlobalFlags() {
	globalFlags := c.root.PersistentFlags()
	globalFlags.StringVar(&c.cfgFile, "config", "", "config file (default is $HOME/.zkpass-node.yaml)")
}

func (c *command) setHomeDir() (err error) {
	if c.homeDir != "" {
		return
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.homeDir = dir
	return nil
}

func (c *command) setAllFlags(cmd *cobra.Command) {
	cmd.Flags().String(optionNameDataDir, filepath.Join(c.homeDir, ".zkpass-node"), "data directory")
	cmd.Flags().Uint32(optionNameSessionMax, 32, "max session count")
	cmd.Flags().Uint32(optionNameSessionTimeout, 120, "session timeout in seconds")
	cmd.Flags().Uint32(optionNameSessionLife, 300, "session life in seconds")
}
