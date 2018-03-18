package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
    Use:   "status [target name]",
    Short: "Display database migrations status for given target name",
    Args:  checkTarget,
    RunE: func(cmd *cobra.Command, args []string) error {
        trg, err := getTarget(args[0])
        if err != nil {
            return err
        }
        for _, st := range trg.Status() {
            msg := "Version: %d, applied: %s\n"
            applied := st.AppliedAt().String()
            if st.AppliedAt().IsZero() {
                applied = "No"
            }
            fmt.Printf(msg, st.Version(), applied)
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
