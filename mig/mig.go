package mig

import (
    "bytes"
    "io/ioutil"
    "path"
    "text/template"
    "time"

    "github.com/pkg/errors"
)

// A Mig is a migrations manager.
type Mig struct {
    cfg    Configurer
    prv    *DriverProvider
    migDir string
}

// NewMig returns new Mig instance.
func NewMig(cfg Configurer) (*Mig, error) {
    m := &Mig{
        cfg:    cfg,
        prv:    NewDriverProvider(cfg),
        migDir: path.Join(cfg.MigDir(), "migrations"),
    }
    // Sort all registered migrations.
    registry.sort()
    return m, nil
}

// NewMigFromConfig instantiates new Mig based on provided config path.
func NewMigFromConfig(path string) (*Mig, error) {
    cfg, err := NewConfigLoader().Load(path)
    if err != nil {
        return nil, err
    }
    return NewMig(cfg)
}

// CreateMigration creates new migration file for given target name.
func (m *Mig) CreateMigration(target string) error {
    drv, err := m.prv.Driver(target)
    if err != nil {
        return err
    }
    if err := m.ensureFs(); err != nil {
        return err
    }
    version := time.Now().UnixNano()
    if err := drv.Creator().CreateMigration(version); err != nil {
        return err
    }
    m.createMain()
    return nil
}

func (m *Mig) Initialize(target string) error {
    drv, err := m.prv.Driver(target)
    if err != nil {
        return err
    }
    if err := m.ensureDb(drv); err != nil {
        return err
    }
    return drv.Initialize()
}

func (m *Mig) Migrate(target string) error {
    drv, err := m.prv.Driver(target)
    if err != nil {
        return err
    }
    if err := m.ensureDb(drv); err != nil {
        return err
    }
    rows, err := drv.Applied()
    if err != nil {
        return err
    }
    if err := registry.applyFromDb(drv, target, rows); err != nil {
        return err
    }
    if err := registry.applyAll(target); err != nil {
        return err
    }
    return nil
}

// ensureFs ensures directory structure is ready for migration files.
func (m *Mig) ensureFs() error {
    if err := checkCreateDir(m.migDir); err != nil {
        return err
    }
    return nil
}

// ensureDb ensures database driver is initialized for migrations.
func (m *Mig) ensureDb(drv Driver) error {
    if err := drv.Open(); err != nil {
        return err
    }
    return nil
}

// createMain creates main.go.
func (m *Mig) createMain() error {
    main := path.Join(m.migDir, "main.go")
    var data = struct {
        Names []string
    }{}
    for _, n := range m.cfg.Targets() {
        if ok, _ := isDir(path.Join(m.migDir, n)); ok {
            data.Names = append(data.Names, n)
        }
    }
    var buf bytes.Buffer
    if err := mainTpl.Execute(&buf, data); err != nil {
        return errors.WithStack(err)
    }
    if err := ioutil.WriteFile(main, buf.Bytes(), 0666); err != nil {
        return errors.WithStack(err)
    }
    return nil
}

var mainTpl = template.Must(template.New("registry-mysqlDriver-struct-tpl").Parse(`package main

import ({{ range $name := .Names }}
    _ "./{{ $name }}"{{ end }}
    "github.com/rzajac/mig/cmd/mig/cmd"
)

// ======================= DO NOT EDIT THIS FILE =======================

func main() {
    cmd.Execute()
}
`))
