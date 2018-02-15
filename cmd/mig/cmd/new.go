package cmd

import (
    "fmt"

    "github.com/pkg/errors"
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Add new migration",
    Long:  `Add new migration.`,
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("requires one arg")
        }
        if args[0] != "mysql" {
            return errors.New("currently only mysql is supported")
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(args)
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
