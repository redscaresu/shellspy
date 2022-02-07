package shellspy_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"shellspy"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCommandFromString(t *testing.T) {

	input := "echo hello world\n"
	want := "/bin/echo hello world"
	got := shellspy.CommandFromString(input)

	if want != got.String() {
		t.Fatal("something gone wrong")
	}

}

func TestRunCommand(t *testing.T) {

	cmd := exec.Command("echo", "hello world")
	want := "hello world\n"
	got, _ := shellspy.RunFromCmd(cmd)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWriteShellScript(t *testing.T) {

	wantBuf := &bytes.Buffer{}
	gotBuf := &bytes.Buffer{}
	wantBuf.WriteString("hello world\n")
	session := shellspy.NewSession()

	session.Input = strings.NewReader("echo hello world")
	session.File = shellspy.CreateFile()
	session.Run()
	os.Remove(session.File.Name())
	fmt.Fprint(gotBuf, session.TranscriptOutput)

	want := wantBuf.String()
	got := gotBuf.String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
