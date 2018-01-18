package mig

import (
    "fmt"
    //"github.com/rzajac/mig/cmd/mig/migration"
)

// Supported dialects
var dialects = [...]string{"mysql"}

// A Mig is a migrations manager.
type Mig struct {
    dir     *dir
    dialect string
}

// NewMig creates new Mig instance.
func NewMig(migDir, dialect string) (*Mig, error) {
    var err error
    if IsSupDialect(dialect) == false {
        return nil, fmt.Errorf("unsupported dialect: %s", dialect)
    }
    m := &Mig{dialect: dialect}
    if m.dir, err = newDir(migDir); err != nil {
        return nil, err
    }
    return m, nil
}

// Init initializes migration directory.
func (m *Mig) Initialize() error {
    return m.dir.init(m.dialect)
}

// New creates new migration file for given dialect.
func (m *Mig) New() (string, error) {
    return m.dir.newMigration(m.dialect)
}

func (m *Mig) Migrate() error {
    //s := reflect.ValueOf(&migration.MigMySQL{})
    //
    //fmt.Println(s.NumMethod())
    //for i := 0; i < s.NumMethod(); i++ {
    //    metName := s.Type().Method(i).Name
    //    fmt.Println(metName)
    //
    //    v := s.Method(i).Call([]reflect.Value{})
    //    isErr := !v[0].IsNil()
    //    fmt.Printf("is error: %v\n", isErr)
    //
    //    if isErr {
    //        err, _ := v[0].Interface().(error)
    //        fmt.Printf("%v\n", err)
    //    }
    //    fmt.Println()
    //}
    return nil
}
