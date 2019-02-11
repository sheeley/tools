package bear

import (
	"net/url"
	"os/exec"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

const bearPrefix = "bear://x-callback-url/"

func Open(action string, parameters map[string]string) error {
	u, err := url.Parse(bearPrefix + strings.TrimPrefix(action, "/"))

	if err != nil {
		return errs.Wrap(err)
	}

	if len(parameters) > 0 {
		q := u.Query()
		for k, v := range parameters {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	u.RawQuery = strings.Replace(u.RawQuery, "+", "%20", -1)

	// fmt.Println(u.String())
	_, err = exec.Command("open", u.String()).CombinedOutput()
	if err != nil {
		return errs.Wrap(err)
	}

	// fmt.Println(string(out))
	return nil
}
