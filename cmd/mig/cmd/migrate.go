package cmd

import (
    "github.com/pkg/errors"
    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate [target name]",
    Short: "Migrate target by name",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("requires target name argument")
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        m, err := mig.NewMigFromConfig(viper.ConfigFileUsed())
        if err != nil {
            printErr(err)
            return
        }
        if err := m.Migrate(args[0]); err != nil {
            printErr(err)
            return
        }
        return
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
