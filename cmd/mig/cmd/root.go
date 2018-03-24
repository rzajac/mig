package cmd

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/pkg/errors"
    "github.com/rzajac/mig/mig"
    "github.com/rzajac/mig/version"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// cfgFile holds the path to configuration file.
var cfgFile string

func init() {
    cobra.OnInitialize(loadConfig)
    rootCmd.SetVersionTemplate(`{{.Version}}`)
    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (default is ./mig.yaml)")
    rootCmd.PersistentFlags().BoolP("version", "v", false, "version")
    rootCmd.PersistentFlags().BoolP("debug", "d", false, "run nin debug mode")
    cfgFile = os.Getenv("MIG_CONFIG")
}

// rootCmd is the main command for the box-auth binary.
var rootCmd = &cobra.Command{
    Use:     "mig",
    Version: getVersion(),
    Short:   "database migration tool",
    Long:    "Database migration tool.",
}

// Execute executes root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

// getTarget returns migration target by name.
func getTarget(name string) (mig.Target, error) {
    m, err := mig.NewMig(cfgFile)
    if err != nil {
        return nil, err
    }
    trg, err := m.Target(name)
    if err != nil {
        return nil, err
    }
    return trg, nil
}

// loadConfig reads in config file.
func loadConfig() {
    viper.SetConfigName("mig")
    viper.AddConfigPath(".")
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    }
    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err != nil {
        log.Fatal(err)
    }
    cfgFile = viper.ConfigFileUsed()
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

// stackTracer is an interface implemented by errors with stack trace.
type stackTracer interface {
    StackTrace() errors.StackTrace
}

func printErr(err error) {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    if err, ok := err.(stackTracer); ok {
        fmt.Fprintf(os.Stderr, "%+v\n", err)
    }
}

// checkTarget checks if the first parameter is a target name.
func checkTarget(_ *cobra.Command, args []string) error {
    if len(args) < 1 {
        return errors.New("requires target name argument")
    }
    return nil
}

func getArgsVersion(args []string, idx int, def int64) (int64, error) {
    if len(args)-1 < idx {
        return def, nil
    }
    return strconv.ParseInt(args[idx], 10, 64)
}
