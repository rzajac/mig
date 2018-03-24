package mig

import "github.com/spf13/afero"

// prepareFs creates in memory filesystem with mig.yaml file.
func prepareFs(configFile string) {
    // Set global filesystem abstraction to memory backed.
    fs = &afero.Afero{Fs: afero.NewMemMapFs()}
    fs.MkdirAll("a/b", 0755)
    afero.WriteFile(fs, "a/b/mig.yaml", []byte(configFile), 0644)
}

var tstYAML = `
dir: migrations
targets:
  t1:
    dialect: mysql
    dsn: t1_dsn
  t2:
    dialect: bad
    dsn: t1_dsn
`
