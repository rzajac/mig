package main

import (
    "fmt"
    "os"

    "github.com/codegangsta/cli"
    "github.com/rzajac/mig/cmd/mig/migration"
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
            Name:        "new",
            Usage:       "Add new migration",
            Description: "adds new migration file",
            Action:      NewMigration,
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "dir, d",
                    Usage: "migrations directory",
                },
            },
        },
    }

    app.Run(os.Args)
}

func NewMigration(ctx *cli.Context) error {
    m, err := mig.NewMigrationFile(ctx.String("dir"))
    if err != nil {
        fmt.Printf("%v\n", err)
        return err
    }
    m.Create()

    ss := &migration.Migration{}
    ss.Mig1516144287290322398()

    return nil
}
