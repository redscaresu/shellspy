package shellspy_test

import (
	"bytes"
	"os/exec"
	"shellspy"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	// var want bytes.Buffer
	// want.WriteString("hello world\n")
	// writer := &bytes.Buffer{}

	// io.WriteString(writer, "hello world\n")
	cmd := exec.Command("echo", "hello world")
	want := "hello world\n"
	got, _ := shellspy.RunFromCmd(cmd)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWriteShellScript(t *testing.T) {

	want := "echo hello world\nhello world"

	writer := &bytes.Buffer{}
	twriter := &bytes.Buffer{}

	session := shellspy.NewSession()
	session.Input = strings.NewReader("echo hello world")
	session.Output = writer
	session.TranscriptOutput = twriter
	session.Run()

	got := twriter.String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
