package cmd

import (
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate [target] [version]",
    Short: "Migrate target",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        trg, err := getTarget(args[0])
        if err != nil {
            return err
        }
        version, err := getArgsVersion(args, 1, 0)
        if err != nil {
            return err
        }
        if err := trg.Migrate(version); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
