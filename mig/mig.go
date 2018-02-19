package mig

// Supported dialects
var dialects = [...]string{"mysql"}

// A Mig is a migrations manager.
type Mig struct {
    config *config
}

// NewMig creates new Mig instance.
func NewMig(config *config) (*Mig, error) {
    m := &Mig{
        config: config,
    }
    return m, nil
}

// Init initializes migration directory.
func (m *Mig) Initialize() error {
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
