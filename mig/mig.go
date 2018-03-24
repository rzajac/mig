package mig

import (
    "path"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
    "gopkg.in/yaml.v2"
)

// File system abstraction.
var fs = &afero.Afero{afero.NewOsFs()}

// mig represents YAML configuration file.
type mig struct {
    Dir string `yaml:"dir"` // Absolute migrations path.
    Targets map[string]*struct {
        Dialect string `yaml:"dialect"` // Database dialect.
        Dsn     string `yaml:"dsn"`     // Database connection string.
    } `yaml:"targets"`      // List of defined targets.
}

// NewMig loads YAML configuration.
func NewMig(configPath string) (*mig, error) {
    content, err := afero.ReadFile(fs, configPath)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    cfg := &mig{}
    if err = yaml.UnmarshalStrict(content, cfg); err != nil {
        return nil, errors.WithStack(err)
    }
    // In config file we only have directory name.
    // This builds absolute path to migrations directory
    // based on configuration file path.
    cfg.Dir = path.Join(path.Dir(configPath), cfg.Dir)
    return cfg, nil
}

func (m *mig) MigDir() string {
    return m.Dir
}

func (m *mig) Target(trgName string) (Target, error) {
    trgCfg, ok := m.Targets[trgName]
    if !ok {
        return nil, ErrUnknownTarget
    }
    constructor, err := GetDriver(trgCfg.Dialect)
    if err != nil {
        return nil, err
    }
    drv := constructor(trgName, trgCfg.Dsn)
    return NewTarget(path.Join(m.Dir, trgName), drv, GetMigrations(trgName))
}

// Names returns all defined target names.
func (m *mig) Names() []string {
    var names []string
    for name := range m.Targets {
        names = append(names, name)
    }
    return names
}
