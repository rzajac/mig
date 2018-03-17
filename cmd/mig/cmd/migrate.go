package cmd

import (
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate [target]",
    Short: "Migrate target",
    Args:  checkTarget,
    Run: func(cmd *cobra.Command, args []string) {
        m, err := NewMigFromConfig(cfgFile, args[0])
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
