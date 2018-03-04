package cmd

import (
    "github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
    Use:   "status [target name]",
    Short: "Display database migrations status for given target name",
    Args: func(cmd *cobra.Command, args []string) error {
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        //m, err := mig.NewMigFromConfig(viper.ConfigFileUsed())
        //if err != nil {
        //    return err
        //}
        return nil
    },
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
