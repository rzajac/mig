package cmd

import "github.com/spf13/cobra"

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migrate database to latest version",
    Run: func(cmd *cobra.Command, args []string) {

    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
