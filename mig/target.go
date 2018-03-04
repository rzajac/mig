package mig

// target represents migration target configuration.
type target struct {
    TDialect string `yaml:"dialect"` // Target database dialect.
    TDsn     string `yaml:"dsn"`     // Target database connection string.
    name     string                  // Target name.
    migDir   string                  // Absolute path to migrations directory.
}

func (t *target) Dialect() string {
    return t.TDialect
}

func (t *target) Dsn() string {
    return t.TDsn
}

func (t *target) Name() string {
    return t.name
}

func (t *target) MigDir() string {
    return t.migDir
}
