package mig

import (
    "go/build"
    "io/ioutil"
    "os"
    "path"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)

// config represents migration configuration.
type config struct {
    Path      string                                  // Absolute path to configuration file.
    Package   string               `yaml:"package"`   // Migrations root directory.
    Databases map[string]*DBConfig `yaml:"databases"` // List of database configurations.
}

// DBConfig represents database configuration.
type DBConfig struct {
    Dialect string `yaml:"dialect"`
    Host    string `yaml:"host"`
    User    string `yaml:"user"`
    Pass    string `yaml:"pass"`
    Name    string `yaml:"name"`
}

// migDir returns absolute path to the directory for
// migrations for given database name.
func (cfg *config) migDir(name string) (string, error) {
    db := cfg.getDBConfig(name)
    if db == nil {
        return "", errors.Errorf("unknown database %s", name)
    }
    return path.Join(cfg.Path, db.Name, db.Dialect), nil
}

// getDBConfig returns database configuration by name.
func (cfg *config) getDBConfig(name string) *DBConfig {
    if db, ok := cfg.Databases[name]; ok {
        return db
    }
    return nil
}

// validate validates mig configuration.
func (cfg *config) validate() error {
    for _, db := range cfg.Databases {
        if err := db.validate(); err != nil {
            return err
        }
    }
    // TODO: validate Config.
    return nil
}

// validate validates database configuration.
func (cfg *DBConfig) validate() error {
    return nil
}

// LoadConfig loads mig configuration file located at cfgPath.
func LoadConfig(cfgPath string) (*config, error) {
    file, err := ioutil.ReadFile(cfgPath)
    if err != nil {
        return nil, err
    }
    cfg := &config{}
    if err := yaml.UnmarshalStrict(file, cfg); err != nil {
        return nil, err
    }
    gp := os.Getenv("GOPATH")
    if gp == "" {
        gp = build.Default.GOPATH
    }
    cfg.Path = path.Join(gp, "src", cfg.Package)
    if err := cfg.validate(); err != nil {
        return nil, err
    }
    return cfg, nil
}
