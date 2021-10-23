package shellspy_test

import (
	"shellspy"
	"testing"
)

func TestCommandFromString(t *testing.T) {
	cmd := "hello world"
	got, err := shellspy.CommandFromString(cmd)
	if err != nil {
		t.Fatal()
	}

	want, err := shellspy.CommandFromString("hello world")
	if err != nil {
		t.Fatal()
	}

	if want != got {
		t.Fatal("want not equal to got")
	}

}
