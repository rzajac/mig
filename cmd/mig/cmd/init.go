package cmd

import (
    "github.com/davecgh/go-spew/spew"
    "github.com/rzajac/mig/mig"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize migrations",
    Long:  `Initialize migrations.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg, err := mig.LoadConfig(viper.ConfigFileUsed())
        if err != nil {
            return err
        }
        dir, _ := cfg.MigDir("auth")
        spew.Dump(cfg, err, dir)
        //for _, cfg := range cfg.Databases {
        //    migPath := path.Join(path.Dir(cfgPath), out.Path, cfg.Dialect, cfg.Name)
        //    m, err := mig.NewMig(migPath, cfg.Dialect)
        //    if err != nil {
        //        return err
        //    }
        //    if err := m.Initialize(); err != nil {
        //        return err
        //    }
        //}
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
