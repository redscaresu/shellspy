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

		want := fmt.Sprint(exec.Command(v))
		fromApp, _ := shellspy.CommandFromString(v)
		got := fmt.Sprint(fromApp)

		if want != got {
			t.Fatal("want not equal to got")
		}
	}
}
