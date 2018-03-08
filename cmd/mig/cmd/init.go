package cmd

import (
    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var initCmd = &cobra.Command{
    Use:   "init [target name]",
    Short: "Initialize target database",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        m, err := mig.NewMigFromConfig(viper.ConfigFileUsed())
        if err != nil {
            return err
        }
        if err := m.Initialize(args[0]); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
