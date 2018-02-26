package mig

type mysqlDriver struct {
    cfg DbConfigurator
}

func NewMYSQLDriver(config DbConfigurator) *mysqlDriver {
    return &mysqlDriver{config}
}

func (m *mysqlDriver) Open() error {
    return nil
}

func (m *mysqlDriver) Close() error {
    return nil
}

func (m *mysqlDriver) Version() (int64, error) {
    return 0, nil
}

func (m *mysqlDriver) Apply(migration Migration) error {
    return nil
}

func (m *mysqlDriver) Initialize() error {
    return nil
}

func (m *mysqlDriver) Creator() Creator {
    return NewMySQLCreator(m.cfg)
}
