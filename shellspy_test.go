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

	want := "hello world"
	cmd := exec.Command("echo hello world")
	got, stdErr := shellspy.RunFromCmd(cmd)
	if stdErr != "<nil>" {
		t.Fatal("something gone wrong")
	}
	if want != got {
		t.Fatal("something gone wrong")
	}

}
