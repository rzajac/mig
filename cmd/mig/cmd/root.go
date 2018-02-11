package cmd

import (
    "encoding/json"
    "os"
    "strings"

    "github.com/rzajac/mig/version"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// cfgFile holds the path to configuration file.
var cfgFile string

// log provides global logger.
var log *logrus.Entry

func init() {
    logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.999999-07:00"})
    log = logrus.WithFields(logrus.Fields{"service": "mig"})

    cobra.OnInitialize(initConfig)
    rootCmd.SetVersionTemplate(`{{.Version}}`)
    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (default is ./mig.yaml)")
    rootCmd.PersistentFlags().BoolP("version", "v", false, "version")
    rootCmd.PersistentFlags().BoolP("debug", "d", false, "run nin debug mode")
    viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

// rootCmd is the main command for the box-auth binary.
var rootCmd = &cobra.Command{
    Use:     "mig",
    Version: getVersion(),
    Short:   "database migration tool",
    Long:    `Database migration tool.`,
}

// Execute executes root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Error(err)
        os.Exit(1)
    }
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    // Add a prefix while reading from the environment variables.
    viper.SetEnvPrefix("MIG")
    // Replace dot with underscore when looking for environmental variables.
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    // Name of the configuration file and where to look for it.
    viper.SetConfigName("mig")
    viper.AddConfigPath("/usr/etc/mig")
    viper.AddConfigPath(".")
    // Load config file if it was explicitly passed.
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    }
    viper.AutomaticEnv()

    // If a config file is found, read it in.
    viper.ReadInConfig()
}

// getVersion returns JSON formatted application version.
func getVersion() string {
    v := struct {
        Version, BuildDate, GitHash, GitTreeState string
    }{
        version.Version,
        version.BuildDate,
        version.GitHash,
        version.GitTreeState,
    }

    j, _ := json.Marshal(v)

    return string(j)
}
