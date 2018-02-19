package mig

import (
    "io/ioutil"
    "os"
    "path"

    "github.com/mitchellh/go-homedir"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)

// config represents mig configuration.
type config struct {
    // Absolute path to configuration file.
    Path string
    // Migrations root directory.
    Package string `yaml:"package"`
    // List of database configurations.
    Databases map[string]*dbConfig `yaml:"databases"`
}

// dbConfig represents database configuration.
type dbConfig struct {
    Dialect string `yaml:"dialect"`
    Host    string `yaml:"host"`
    User    string `yaml:"user"`
    Pass    string `yaml:"pass"`
    Name    string `yaml:"name"`
    migDir  string
    Pack    string
}

// newConfig initializes and validates mig configuration.
// Takes absolute path to configuration file.
func newConfig(path string) (*config, error) {
    cfg, err := loadConfig(path)
    if err != nil {
        return nil, err
    }
    if err := cfg.init(); err != nil {
        return nil, err
    }
    return cfg, nil
}

// init initializes and validates mig configuration.
func (cfg *config) init() (err error) {
    gp := os.Getenv("GOPATH")
    if gp == "" {
        gp, err = homedir.Dir()
        if err != nil {
            return err
        }
        gp = path.Join(gp, "go")
    }
    cfg.Path = path.Join(gp, "src", cfg.Package)
    for name, db := range cfg.Databases {
        db.migDir = path.Join(cfg.Path, "migrations", name)
        db.Pack = path.Dir(db.migDir)
        if !IsSupDialect(db.Dialect) {
            return errors.Errorf("unsupported dialect %s", db.Dialect)
        }
    }
    return nil
}

// getDBConfig returns database configuration by name.
func (cfg *config) getDBConfig(name string) *dbConfig {
    if db, ok := cfg.Databases[name]; ok {
        return db
    }
    return nil
}

// loadConfig loads and validates mig configuration file.
func loadConfig(cfgPath string) (*config, error) {
    file, err := ioutil.ReadFile(cfgPath)
    if err != nil {
        return nil, err
    }
    cfg := &config{}
    if err := yaml.UnmarshalStrict(file, cfg); err != nil {
        return nil, err
    }

    if err := cfg.init(); err != nil {
        return nil, err
    }
    return cfg, nil
}
