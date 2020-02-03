package input

import (
	"bufio"
	"os"

	"github.com/richardwilkes/toolbox/errs"
)

func ReadChar() (rune, error) {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return char, errs.Wrap(err)
	}
	return char, nil
}
