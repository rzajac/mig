package cmd

import (
    "github.com/pkg/errors"
    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var newCmd = &cobra.Command{
    Use:   "new [name]",
    Short: "Create new migration",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("requires connection name argument")
        }
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        m, err := mig.NewMig(viper.ConfigFileUsed())
        if err != nil {
            return err
        }
        if err := m.Initialize(); err != nil {
            return err
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
