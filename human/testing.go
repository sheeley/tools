package human

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func AssertOrDiff(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(s1, s2, false)
		// fmt.Printf("\n=====\n%s\n=====\n%s\n=====\n", s1, s2)

		fmt.Println(dmp.DiffPrettyText(diffs))
		assert.Equal(t, s1, s2)

	}
}
