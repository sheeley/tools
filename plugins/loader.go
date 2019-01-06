package plugins

import (
	"plugin"

	"github.com/richardwilkes/toolbox/errs"
)

type Loader struct {
	Path string
}

func (l *Loader) Open() (*plugin.Plugin, error) {
	p, err := plugin.Open(l.Path)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return p, nil
}
