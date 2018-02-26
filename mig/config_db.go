package mig

// dbConfig represents a database configuration.
type dbConfig struct {
    DbDialect string `yaml:"dialect"` // Database dialect.
    DbDsn     string `yaml:"dsn"`     // Database connection string.
    name      string                  // The configuration name.
    migDir    string                  // Absolute path to migrations directory.
}

// Dialect implements DbConfigurator interface.
func (c *dbConfig) Dialect() string {
    return c.DbDialect
}

// Dsn implements DbConfigurator interface.
func (c *dbConfig) Dsn() string {
    return c.DbDsn
}

// Name implements DbConfigurator interface.
func (c *dbConfig) Name() string {
    return c.name
}

// MigDir implements DbConfigurator interface.
func (c *dbConfig) MigDir() string {
    return c.migDir
}
