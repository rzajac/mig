package mig

import (
    "bytes"
    "path"
    "text/template"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
)

// checkCreateDir creates directory if doesn't exist.
func checkCreateDir(fs afero.Fs, path string) error {
    ok, err := afero.DirExists(fs, path)
    if err != nil {
        return err
    }
    if !ok {
        if err := fs.MkdirAll(path, 0777); err != nil {
            return errors.WithStack(err)
        }
    }
    return nil
}

// createMain creates main.go.
func createMain(fs afero.Fs, dir string, names ...string) error {
    var data = struct {
        Names []string
    }{}
    data.Names = append(data.Names, names...)
    for n := range migrations {
        data.Names = append(data.Names, n)
    }
    var buf bytes.Buffer
    if err := mainTpl.Execute(&buf, data); err != nil {
        return errors.WithStack(err)
    }
    main := path.Join(dir, "main.go")
    if err := afero.WriteFile(fs, main, buf.Bytes(), 0666); err != nil {
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
