package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/codegangsta/cli"
    "github.com/pkg/errors"
    "github.com/rzajac/mig/mig"
)

func main() {
    app := cli.NewApp()
    app.Name = "mig"
    app.Usage = "database migration tool"
    app.Version = "0.0.1"
    app.Authors = []cli.Author{
        {Name: "Rafal Zajac", Email: "rzajac@gmail.com"},
    }
    app.Copyright = "(c) 2018 Rafal Zajac <rzajac@gmail.com>"

    app.Commands = []cli.Command{
        {
            Name:        "init",
            Usage:       "Initialize migrations.",
            ArgsUsage:   "dialect dir",
            Description: "Initializes migrations directory for given dialect.",
            Action:      InitMigrationsCmd,
        },
        {
            Name:        "new",
            Usage:       "Add new migration",
            ArgsUsage:   "dialect dir",
            Description: "adds new migration file for given dialect",
            Action:      NewMigrationCmd,
        },
        {
            Name:        "migrate",
            Usage:       "Apply migrations",
            ArgsUsage:   "dialect dir",
            Description: "applies not applied migrations",
            Action:      Migrate,
        },
    }

    if err := app.Run(os.Args); err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }
}

// getMig returns Mig instance.
func getMig(ctx *cli.Context) (*mig.Mig, error) {
    dialect := ctx.Args().First()
    dir := ctx.Args().Get(1)
    if err := validateArgs(dialect, dir); err != nil {
        return nil, err
    }
    m, err := mig.NewMig(dir, dialect)
    if err != nil {
        return nil, err
    }
    return m, nil
}

// validateArgs validate command arguments.
func validateArgs(dialect, dir string) error {
    if !mig.IsSupDialect(dialect) {
        return fmt.Errorf("unsupported dialect: %s", dialect)
    }
    if dir == "" {
        return errors.New("directory argument must be provided")
    }
    return nil
}

// InitMigrationsCmd initialize migration directory for given dialect.
func InitMigrationsCmd(ctx *cli.Context) error {
    m, err := getMig(ctx)
    if err != nil {
        return err
    }
    return m.Initialize()
}

// NewMigrationCmd creates new migration file.
func NewMigrationCmd(ctx *cli.Context) error {
    m, err := getMig(ctx)
    if err != nil {
        return err
    }
    file, err := m.New()
    if err != nil {
        return err
    }

    tmp, _ := os.Getwd()
    tmp, _ = filepath.Rel(tmp, file)
    fmt.Printf("created migration file %s\n", tmp)
    return nil
}

func Migrate(ctx *cli.Context) error {
    m, err := getMig(ctx)
    if err != nil {
        return err
    }
    return m.Migrate()
}
