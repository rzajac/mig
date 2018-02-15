package cmd

import "github.com/spf13/cobra"

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migrate database",
    Long:  `Apply not applied migrations.`,
    Run: func(cmd *cobra.Command, args []string) {

    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
