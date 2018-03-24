package cmd

import (
    "os"
    "strconv"

    "github.com/olekukonko/tablewriter"
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

        status, err := trg.Status()
        if err != nil {
            return err
        }

        data := make([][]string, 0)
        for _, st := range status {
            applied := st.AppliedAt().String()
            if st.AppliedAt().IsZero() {
                applied = "No"
            }
            row := []string{strconv.FormatInt(st.Version(), 10), st.Info(), applied}
            data = append(data, row)
        }

        table := tablewriter.NewWriter(os.Stdout)
        table.SetHeader([]string{"Version", "Description", "Applied"})
        table.SetBorder(true)
        table.AppendBulk(data)
        table.Render()

        return nil
    },
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
