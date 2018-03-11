package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate [target name]",
    Short: "Migrate target by name",
    Args:  checkTarget,
    Run: func(cmd *cobra.Command, args []string) {
        m, err := NewMigFromConfig(viper.ConfigFileUsed(), args[0])
        if err != nil {
            printErr(err)
            return
        }
        if err := m.Migrate(0); err != nil {
            printErr(err)
            return
        }
        return
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
