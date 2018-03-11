package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var newCmd = &cobra.Command{
    Use:   "new [target name]",
    Short: "Create new migration for given target name",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        m, err := NewMigFromConfig(viper.ConfigFileUsed(), args[0])
        if err != nil {
            return err
        }
        if err := m.CreateMigration(); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
