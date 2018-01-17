package mig

// A Mig is a migrations manager.
type Mig struct {
    dir     *Dir
    dialect string
}

// NewMig creates new Mig instance.
func NewMig(dir, dialect string) (*Mig, error) {
    d, err := NewDir(dir)
    if err != nil {
        return nil, err
    }
    return &Mig{d, dialect}, nil
}

func (m *Mig) Initialize() error {
    return m.dir.Initialize(m.dialect)
}

func (m *Mig) New() error {
    return m.dir.NewMigration(m.dialect)
}