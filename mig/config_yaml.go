package mig

import (
    "io/ioutil"
    "path"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)

// NewYAMLConfigurator returns instance of YAML Configurator.
func NewYAMLConfigurator(url string) (Configurator, error) {
    var err error
    var data []byte
    cfg := &yamlCfg{}
    if data, err = ioutil.ReadFile(url); err != nil {
        return nil, errors.WithStack(err)
    }
    if err = yaml.UnmarshalStrict(data, cfg); err != nil {
        return nil, errors.WithStack(err)
    }
    cfg.baseDir = path.Dir(url)
    for n, db := range cfg.DBs {
        db.name = n
        db.migDir = path.Join(cfg.baseDir, "migrations", n)
    }
    return cfg, nil
}

// yamlCfg represents YAML configuration file.
type yamlCfg struct {
    baseDir string
    DBs     map[string]*dbConfig `yaml:"databases"`
}

// BaseDir implements Configurator interface.
func (c *yamlCfg) BaseDir() string {
    return c.baseDir
}

// DbConfig implements Configurator interface.
func (c *yamlCfg) DbConfig(name string) (DbConfigurator, error) {
    db, ok := c.DBs[name]
    if !ok {
        return nil, errors.Errorf("no database config named %s", name)
    }
    return db, nil
}
