package cmd

import (
    "fmt"
    "strings"

    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
)

var targetsCmd = &cobra.Command{
    Use:   "targets",
    Short: "Display database target names",
    RunE: func(cmd *cobra.Command, args []string) error {
        m, err := mig.NewMig(cfgFile)
        if err != nil {
            return err
        }
        fmt.Printf("%s\n", strings.Join(m.Names(), " "))
        return nil
    },
}

func init() {
    rootCmd.AddCommand(targetsCmd)
}
