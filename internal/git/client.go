package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// returns the raw unified diff of staged changes
func GetStagedDiff() (string, error) {
	// --cached gets staged changes
	// --unified=3 provides 3 lines of context around changes
	cmd := exec.Command("git", "diff", "--cached", "--unified=3")
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
