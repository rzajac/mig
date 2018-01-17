package main

import (
    "os"

    "github.com/codegangsta/cli"
    "github.com/rzajac/mig/mig"
    "fmt"
    "github.com/pkg/errors"
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
            Action:      InitMigrations,
        },
        {
            Name:        "new",
            Usage:       "Add new migration",
            ArgsUsage:   "dialect dir",
            Description: "adds new migration file for given dialect",
            Action:      NewMigration,
        },
    }

    if err := app.Run(os.Args); err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }
}

func getMig(ctx *cli.Context) (*mig.Mig, error) {
    dialect := ctx.Args().First()
    dir := ctx.Args().Get(1)
    if err := heckDialectDir(dialect, dir); err != nil {
        return nil, err
    }
    m, err := mig.NewMig(dir, dialect)
    if err != nil {
        return nil, err
    }
    return m, nil
}

func heckDialectDir(dialect, dir string) error {
    if dialect == "" {
        return errors.New("dialect must be provided")
    }
    if dir == "" {
        return errors.New("directory must be provided")
    }
    return nil
}

func InitMigrations(ctx *cli.Context) error {
    m, err := getMig(ctx)
    if err != nil {
        return err
    }
    return m.Initialize()
}

func NewMigration(ctx *cli.Context) error {
    m, err := getMig(ctx)
    if err != nil {
        return err
    }
    return m.New()
}
