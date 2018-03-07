package mig

import (
    "errors"
    "time"
)

var ErrNoMoreMigrations = errors.New("no more migrations")
var ErrNotInitialized = errors.New("database not initialized")

// MigRows represents map of applied migration versions
// and when they were applied.
type MigRows map[int64]time.Time

// The Configurer interface is implemented by configuration providers.
type Configurer interface {
    // MigDir returns absolute path to "migrations" directory.
    // In another words it returns absolute path to directory
    // where all target directories are being created.
    MigDir() string
    // Target returns migration target by name.
    Target(name string) (Target, error)
    // Targets returns all migration target names.
    Targets() []string
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
    Apply(Migrator) error
    // Revert reverts migration.
    Revert(Migrator) error
    // Applied returns all applied migrations.
    // The list is ordered starting with most recent migration.
    // Returns ErrNotInitialized if database is not prepared for migrations.
    Applied() (MigRows, error)
    // Creator returns migration file creator for given Driver.
    Creator() Creator
}

// The Migrator interface is used to apply or revert given migration.
type Migrator interface {
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

// The Migrator interface provides methods to access
// migrations for given database.
//
// Calls to Next and Previous provide a way
// to move back and forth through migration migs.
type Manager interface {
    // SetCurrent sets current database migration version.
    // Can be called at any time to reset the start point on the
    // migration migs.
    SetCurrent(id int64) error
    // Current returns migration version which Migrator currently points to.
    Current() int64
    // Next returns next migration to apply.
    // Next returns ErrNoMoreMigrations error if
    // there are no more migrations on the migs.
    Next() (Migrator, error)
    // Previous returns migration which preceded the current migration.
    // Previous returns ErrNoMoreMigrations error if the current migration is
    // first on the migrations migs.
    Previous() (Migrator, error)
}
