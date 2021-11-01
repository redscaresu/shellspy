package shellspy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CommandFromString(input string) (*exec.Cmd, error) {
	trim := strings.TrimSuffix(input, "\n")
	name := strings.Fields(trim)
	args := name[1:]
	join := strings.Join(args, "")
	cmd := exec.Command(name[0], join)
	return cmd, nil
}

func RunFromCmd(cmd *exec.Cmd) error {

	err := cmd.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "an error has occured, %v\n", err)
		os.Exit(2)
	}
	return err
}
