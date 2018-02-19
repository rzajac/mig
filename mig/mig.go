package mig

import (
    "os"

    "github.com/pkg/errors"
)

// Supported dialects
var dialects = [...]string{"mysql"}

// A Mig is a migrations manager.
type Mig struct {
    config *config
}

// NewMig creates new Mig instance.
func NewMig(configPath string) (*Mig, error) {
    cfg, err := newConfig(configPath)
    if err != nil {
        return nil, err
    }
    return &Mig{cfg}, nil
}

// Init initializes migration directory.
func (m *Mig) Initialize() error {
    for _, db := range m.config.Databases {
        ok, err := IsDir(db.migDir)
        if err != nil {
            return err
        }
        if ok {
            continue
        }
        if err := os.MkdirAll(db.migDir, 0777); err != nil {
            return errors.WithStack(err)
        }
        if err := CreateBaseMig(db.migDir, db.Dialect); err != nil {
            return err
        }
    }
    return nil
}

// New creates new migration file for given dialect.
func (m *Mig) New() (string, error) {
    return "", nil
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
