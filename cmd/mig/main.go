package main

import (
    "github.com/rzajac/mig/cmd/mig/cmd"
)

func main() {
    cmd.Execute()
}

//// getMig returns Mig instance.
//func getMig(ctx *cli.Context) (*mig.Mig, error) {
//    dialect := ctx.Args().First()
//    dir := ctx.Args().Get(1)
//    if err := validateArgs(dialect, dir); err != nil {
//        return nil, err
//    }
//    m, err := mig.NewMig(dir, dialect)
//    if err != nil {
//        return nil, err
//    }
//    return m, nil
//}

//// validateArgs validate command arguments.
//func validateArgs(dialect, dir string) error {
//    if !mig.IsSupDialect(dialect) {
//        return fmt.Errorf("unsupported dialect: %s", dialect)
//    }
//    if dir == "" {
//        return errors.New("directory argument must be provided")
//    }
//    return nil
//}
//
//// InitMigrationsCmd initialize migration directory for given dialect.
//func InitMigrationsCmd(ctx *cli.Context) error {
//    m, err := getMig(ctx)
//    if err != nil {
//        return err
//    }
//    return m.Initialize()
//}
//
//// NewMigrationCmd creates new migration file.
//func NewMigrationCmd(ctx *cli.Context) error {
//    m, err := getMig(ctx)
//    if err != nil {
//        return err
//    }
//    file, err := m.New()
//    if err != nil {
//        return err
//    }
//
//    tmp, _ := os.Getwd()
//    tmp, _ = filepath.Rel(tmp, file)
//    fmt.Printf("created migration file %s\n", tmp)
//    return nil
//}
//
//func Migrate(ctx *cli.Context) error {
//    m, err := getMig(ctx)
//    if err != nil {
//        return err
//    }
//    return m.Migrate()
//}
