package mig

import "time"

// record represents one record in database migrations table.
// Each record is a log of successfully run migration.
type record struct {
    version   int64
    createdAt time.Time
}

func (r *record) Version() int64 {
    return r.version
}

func (r *record) CreatedAt() time.Time {
    return r.createdAt
}
