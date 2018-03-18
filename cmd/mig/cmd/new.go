package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new [target]",
    Short: "Create new migration for given target",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        trg, err := getTarget(args[0])
        if err != nil {
            return err
        }
        if err := trg.CreateMigration(); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
