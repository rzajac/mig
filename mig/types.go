package mig

import (
    "errors"
    "time"
)

var ErrNotInitialized = errors.New("database not initialized")
var ErrUnknownTarget = errors.New("unknown target")

// The Target interface is implemented by target configuration providers.
type Target interface {
    // Name returns target name.
    Name() string
    // TargetDir returns absolute path to a directory where target migrations are.
    TargetDir() string
    // CreateMigration creates migration file.
    CreateMigration() error
    // Migrate migrates target to given version.
    // If version is 0 it will migrate to latest version.
    Migrate(toVersion int64) error
    // Initialize prepares underlying database for migrations.
    Initialize() error
    // Returns migration status for the target.
    Status() []Status
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
    // GenMigration generates migration file.
    GenMigration(version int64) ([]byte, error)
}

type Status interface {
    // Version returns migration version.
    Version() int64
    // AppliedAt returns when migration has been applied.
    // It might return Zero date if the migration has not been applied.
    AppliedAt() time.Time
    // Info returns short (140 characters max) migration description.
    Info() string
}

// The Migration interface is implemented by database migrations.
type Migration interface {
    Status
    // Setup is called by migration manager before calling any other method.
    Setup(driver interface{}, createdAt time.Time)
    // Apply applies migration to driver.
    Apply() error
    // Revert reverts migration.
    Revert() error
}
