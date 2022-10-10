package cmd

import (
	"fmt"
	"zkpass-node"
	"zkpass-node/pkg/node"

	"github.com/spf13/cobra"
)

func (c *command) initStartCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a zkpass-node",
		RunE: func(cmd *cobra.Command, args []string) (e error) {
			if len(args) > 0 {
				return cmd.Help()
			}

			fmt.Println("zkpass-node version", zkpass.Version)

			_, err := node.New(&node.Options{
				DataDir:        c.config.GetString(optionNameDataDir),
				SessionMax:     c.config.GetInt64(optionNameSessionMax),
				SessionTimeout: c.config.GetInt64(optionNameSessionTimeout),
				SessionLife:    c.config.GetInt64(optionNameSessionLife),
			})

			if err != nil {
				return err
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	c.setAllFlags(cmd)
	c.root.AddCommand(cmd)
	return nil
}
