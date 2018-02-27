package mig

import (
    "errors"
    "time"
)

var ErrNoMoreMigrations = errors.New("no more migrations")
var ErrNotInitialized = errors.New("database not initialized")

// The Loader interface is implemented by configuration loaders.
type Loader interface {
    // Load loads configuration.
    Load(url string) (Configurator, error)
}

// The Configurator interface is implemented by configuration providers.
type Configurator interface {
    // BaseDir returns absolute path to base directory.
    BaseDir() string
    // DbConfigurator returns migration configuration by name.
    DbConfig(name string) (DbConfigurator, error)
}

// The DbConfigurator interface is implemented by database configuration
// providers.
type DbConfigurator interface {
    // Name returns configuration name.
    Name() string
    // MigDir returns absolute path migrations directory.
    MigDir() string
    // Dialect returns database dialect.
    Dialect() string
    // Dsn returns database connection string.
    Dsn() string
}

// The Migrator interface provides methods to access
// migrations for given database.
//
// Calls to Next and Previous provide a way
// to move back and forth through migration list.
type Migrator interface {
    // SetCurrent sets current database migration version.
    // Can be called at any time to reset the start point on the
    // migration list.
    SetCurrent(id int64) error
    // Current returns migration version which Migrator currently points to.
    Current() int64
    // Next returns next migration to apply.
    // Next returns ErrNoMoreMigrations error if
    // there are no more migrations on the list.
    Next() (Migration, error)
    // Previous returns migration which preceded the current migration.
    // Previous returns ErrNoMoreMigrations error if the current migration is
    // first on the migrations list.
    Previous() (Migration, error)
}

// The Migration interface describes a migration.
type Migration interface {
    // Setup is called by migration manager before calling any other method.
    Setup(driver interface{}, applied bool)
    // Version returns migration version.
    Version() int64
    // Apply applies migration.
    // Apply will set database migration version to value returned by Version().
    Apply() error
    // Revert reverts migration.
    // Revert will remove database migration version returned by Version().
    Revert() error
    // IsApplied returns true if the migration has been applied.
    IsApplied() bool
    // Description returns short (255 characters max) migration description.
    Description() string
}

// DriverProvider provides database drivers.
type DriverProvider interface {
    // Driver returns Driver for given DbConfigurator.
    Driver(name string) (Driver, error)
}

// Driver represents database driver.
type Driver interface {
    // Open opens connection to database.
    Open() error
    // Close closes database connection.
    Close() error
    // Version returns the current database schema version.
    // May return ErrNotInitialized if the underlying database is not
    // prepared for migrations.
    Version() (int64, error)
    // Apply applies migration to the database.
    Apply(Migration) error
    // Applied returns all applied migrations descriptions.
    Applied() ([]Descriptor, error)
    // Initialize prepares underlying database for migrations.
    Initialize() error
    // Creator returns migration file creator.
    Creator() Creator
}

// The Creator interface is implemented by structures which know how to
// create migration files for given dialect.
type Creator interface {
    // CreateMigration creates migration file.
    // If the migrations directory or base migration file doesn't exist
    // CreateMigration will first create them.
    CreateMigration(version int64) error
}

// Descriptor describes applied migration.
type Descriptor interface {
    // Version returns migration version.
    Version() int
    // Description returns migration description.
    Description() string
    // Created returns when migration has been applied in UTC.
    Created() time.Time
}
