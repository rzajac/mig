package cmd

import (
    "io/ioutil"
    "path"

    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v2"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize migrations",
    Long:  `Initialize migrations.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        cfgPath := viper.ConfigFileUsed()
        cfg, err := ioutil.ReadFile(cfgPath)
        if err != nil {
            return err
        }
        out := &mig.Cfg{}
        if err := yaml.UnmarshalStrict(cfg, out); err != nil {
            return err
        }
        for _, cfg := range out.DBs {
            migPath := path.Join(path.Dir(cfgPath), out.Dir, cfg.Dialect, cfg.Name)
            m, err := mig.NewMig(migPath, cfg.Dialect)
            if err != nil {
                return err
            }
            if err := m.Initialize(); err != nil {
                return err
            }
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
