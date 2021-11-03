package shellspy_test

import (
	"os"
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

func TestWriteShellScript(t *testing.T) {

	file, err := os.Open("testdata/transcript.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	want := file
	got := shellspy.WriteTranscript("hello world/n")

	if *want != got {
		t.Fatal("something gone wrong")
	}

	os.Remove("transcript.txt")

}
