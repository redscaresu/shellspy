package shellspy

import (
	"os/exec"
)

func CommandFromString(input string) (*exec.Cmd, error) {
	cmd := exec.Command(input)
	return cmd, nil
}
