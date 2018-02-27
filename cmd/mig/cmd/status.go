package cmd

import (
    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
    Use:   "status [name]",
    Short: "Display status of database migrations",
    Args: func(cmd *cobra.Command, args []string) error {
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        //m, err := mig.NewMigFromConfig(viper.ConfigFileUsed())
        //if err != nil {
        //    return err
        //}
        mig.List()
        return nil
    },
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
