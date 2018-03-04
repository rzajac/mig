package mig

import "time"

// record represents one record in database migrations table.
// Each record is a log of successfully run migration.
type record struct {
    version   int64
    current   bool
    info      string
    createdAt time.Time
}

func (r *record) Version() int64 {
    return r.version
}

func (r *record) Current() bool {
    return r.current
}

func (r *record) Info() string {
    return r.info
}

func (r *record) CreatedAt() time.Time {
    return r.createdAt
}
