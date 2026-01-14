package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// returns the raw unified diff of staged changes
func GetStagedDiff(contextWidth int, staged bool) (string, error) {
	// this defines the contest width of the git diff extracted
	// for example: --unified=3 provides 3 lines of context around changes
	contextDefinition := fmt.Sprintf("--unified=%d", contextWidth)

	var cmd *exec.Cmd
	if staged {
		cmd = exec.Command("git", "diff", "--cached", contextDefinition)
	} else {
		cmd = exec.Command("git", "diff", contextDefinition)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git error: %s", stderr.String())
	}

	if stdout.Len() == 0 {
		return "", fmt.Errorf("no changes staged. Use 'git add' first")
	}

	return stdout.String(), nil
}
