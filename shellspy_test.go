package shellspy_test

import (
	"os/exec"
	"shellspy"
	"testing"
)

func TestCommandFromString(t *testing.T) {
	got, err := shellspy.CommandFromString("echo hello world")
	if err != nil {
		t.Fatal()
	}

	want := exec.Command("echo", "hello world")

	if want != got {
		t.Fatal("want not equal to got")
	}

}
