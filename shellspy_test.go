package shellspy_test

import (
	"fmt"
	"os/exec"
	"shellspy"
	"testing"
)

func TestCommandFromString(t *testing.T) {

	testStrings := []string{"echo", "cat", "ls -la"}

	for _, v := range testStrings {

		wantString := fmt.Sprint(exec.Command(v))
		got, _ := shellspy.CommandFromString(v)
		gotString := fmt.Sprint(got)

		if wantString != gotString {
			t.Fatal("want not equal to got")
		}
	}
}

func TestRunCommand(t *testing.T) {

	cmd := exec.Command("echo")
	err := shellspy.RunFromCmd(cmd)
	if err != nil {
		t.Fatal("something gone wrong")
	}
}
