package cmd

import (
    "github.com/davecgh/go-spew/spew"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize migrations",
    Long:  `Initialize migrations.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        mig, err := GetMig(viper.ConfigFileUsed())
        if err != nil {
            return err
        }
        spew.Dump(mig, err)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
