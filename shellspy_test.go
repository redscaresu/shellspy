package shellspy_test

import (
	"os/exec"
	"shellspy"
	"testing"
)

func TestCommandFromString(t *testing.T) {

	input := "echo hello world\n"
	want := "/bin/echo hello world"
	got, _ := shellspy.CommandFromString(input)

	if want != got.String() {
		t.Fatal("something gone wrong")

	}

}

func TestRunCommand(t *testing.T) {

	want := "hello world\n"
	cmd := exec.Command("echo", "hello world")
	got, _ := shellspy.RunFromCmd(cmd)

	if want != got {
		t.Fatal("something gone wrong")
	}
}
