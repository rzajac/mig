package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "initialize migrations",
    Long:  `Initialize migrations.`,
    Run: func(cmd *cobra.Command, args []string) {
        dir := viper.GetString("mig.dir")
        fmt.Println(dir)
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
