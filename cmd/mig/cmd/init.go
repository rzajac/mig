package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var initCmd = &cobra.Command{
    Use:   "init [target name]",
    Short: "Initialize target database",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        m, err := NewMigFromConfig(viper.ConfigFileUsed(), args[0])
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
