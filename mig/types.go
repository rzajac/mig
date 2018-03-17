package mig

import (
    "errors"
    "time"
)

var ErrNotInitialized = errors.New("database not initialized")

// The Config interface is implemented by configuration providers.
type Config interface {
    // MigDir returns absolute path to "migrations" directory.
    // In another words it returns absolute path to directory
    // where all target directories are being created.
    MigDir() string
    // Driver returns ready to use database driver for given target name.
    Driver(name string) (Driver, error)
    // TargetNames returns all migration target names.
    TargetNames() []string
}

// The Target interface is implemented by target configuration providers.
type Target interface {
    // Name returns target name.
    Name() string
    // MigDir returns absolute directory path where migrations must be put.
    MigDir() string
    // Dialect returns target database dialect.
    Dialect() string
    // Dsn returns Database Source Name string for the target.
    Dsn() string
}

// Driver represents database driver.
type Driver interface {
    // Open opens connection to database.
    // Maybe called multiple times on the same driver instance.
    Open() error
    // Close closes database connection.
    // May be called multiple times.
    Close() error
    // Initialize prepares underlying database for migrations.
    Initialize() error
    // Apply applies migration to the database.
    Apply(Migration) error
    // Revert reverts migration.
    Revert(Migration) error
    // Merge merges data contained in migration files with data about migration
    // from database.
    Merge([]Migration) error
    // Version returns current database migration version.
    // Returns ErrNotInitialized if database is not prepared for migrations.
    Version() (int64, error)
    // Creator returns migration file creator for given Driver.
    Creator() Creator
}

// The Migration interface is implemented by database migrations.
type Migration interface {
    // Setup is called by migration manager before calling any other method.
    Setup(driver interface{}, createdAt time.Time)
    // Version returns migration version.
    Version() int64
    // AppliedAt returns when migration has been applied.
    // It might return Zero date if the migration has not been applied.
    AppliedAt() time.Time
    // Info returns short (140 characters max) migration description.
    Info() string
    // Apply applies migration to driver.
    Apply() error
    // Revert reverts migration.
    Revert() error
}

// The Creator interface is used to create migration files.
type Creator interface {
    // CreateMigration creates migration file.
    //
    // If the migrations directory or base migration file doesn't exist
    // CreateMigration will first create them.
    CreateMigration(version int64) error
}
