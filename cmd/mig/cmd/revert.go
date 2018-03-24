package cmd

import (
    "github.com/spf13/cobra"
)

var revertCmd = &cobra.Command{
    Use:   "revert target version",
    Short: "Revert target",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        trg, err := getTarget(args[0])
        if err != nil {
            return err
        }
        if err := trg.Migrate(-1); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(revertCmd)
}
