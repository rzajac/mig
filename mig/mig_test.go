package mig

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestMig_NewMig(t *testing.T) {
    // --- Given ---
    prepareFs(tstYAML)

    // --- When ---
    cfg, err := NewMig("a/b/mig.yaml")

    // --- Then ---
    assert.NoError(t, err)
    assert.Exactly(t, "a/b/migrations", cfg.MigDir())
    assert.Exactly(t, 2, len(cfg.Targets))
}

func TestMig_Target_t1(t *testing.T) {
    // --- Given ---
    prepareFs(tstYAML)

    // --- When ---
    cfg, _ := NewMig("a/b/mig.yaml")
    t1, err := cfg.Target("t1")

    // --- Then ---
    assert.NoError(t, err)
    assert.Exactly(t, "t1", t1.Name())
    assert.Exactly(t, "a/b/migrations/t1", t1.TargetDir())
}
