package cmd

import (
    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init [target]",
    Short: "Initialize target database",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        trg, err := getTarget(args[0])
        if err != nil {
            return err
        }
        if err := trg.Initialize(); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
