package git

import (
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

// convert the git output string into structured FileDiffs
func ParseRawDiff(rawDiff string) ([]*gitdiff.File, error) {
	reader := strings.NewReader(rawDiff)

	files, _, err := gitdiff.Parse(reader)
	if err != nil {
		return nil, err
	}

	return files, nil
}
