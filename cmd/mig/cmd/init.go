package cmd

import (
    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init [target]",
    Short: "Initialize target database",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        m, err := NewMigFromConfig(cfgFile, args[0])
        if err != nil {
            return err
        }
        if err := m.Initialize(); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
