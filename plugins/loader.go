package plugins

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

type Loader struct {
	Name, Path string
}

const EnvVar = "PRIVATE_DATA_DIR"

func (l *Loader) CompileAndLoad() (*plugin.Plugin, error) {
	dir, _ := os.LookupEnv(EnvVar)
	d := strings.TrimSpace(dir)
	if d == "" {
		return nil, errs.Newf("Missing '%s' env variable", EnvVar)
	}

	plFile := path.Join(d, l.Name+".so")
	goFilePattern := path.Join(d, l.Name, "*.go")

	goFiles, err := filepath.Glob(goFilePattern)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(goFiles) == 0 {
		return nil, errs.Newf("no go files found: %s", goFilePattern)
	}

	args := append([]string{"build", "-buildmode=plugin", "-o", plFile}, goFiles...)
	cmd := exec.Command("go", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errs.NewWithCause(string(out), err)
	}

	pl, err := plugin.Open(plFile)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return pl, nil
}

// func (l *Loader) Open() (*plugin.Plugin, error) {
// 	p, err := plugin.Open(l.Path)
// 	if err != nil {
// 		return nil, errs.Wrap(err)
// 	}
// 	return p, nil
// }
