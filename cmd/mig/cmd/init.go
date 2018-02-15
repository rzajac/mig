package cmd

import (
    "fmt"
    "io/ioutil"
    "path/filepath"

    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v2"
)

var initCmd = &cobra.Command{
    Use:   "init dialect",
    Short: "Initialize migrations",
    Long:  `Initialize migrations.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        dir, err := filepath.Abs(viper.GetString("dir"))
        if err != nil {
            return err
        }

        m, err := mig.NewMig(dir, viper.GetString("db.users.dialect"))
        if err != nil {
            return err
        }

        cfgPath := viper.ConfigFileUsed()
        fmt.Println(cfgPath)

        cfg, err := ioutil.ReadFile(cfgPath)
        if err != nil {
            return err
        }

        fmt.Println(string(cfg))

        out := &mig.MigConfig{}
        if err := yaml.UnmarshalStrict(cfg, out); err != nil {
            return err
        }

        fmt.Println("Config:")
        fmt.Println(out.DBs["users"].Dialect)

        return m.Initialize()
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
