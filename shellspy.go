package shellspy

import (
	"fmt"
	"os"
	"os/exec"
)

func CommandFromString(input string) (*exec.Cmd, error) {
	cmd := exec.Command(input)
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
