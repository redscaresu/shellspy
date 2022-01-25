package shellspy_test

import (
	"bytes"
	"fmt"
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

	// twriter := &bytes.Buffer{}
	// stdOut := "hello world\n"
	// wantSession := shellspy.NewSession()
	// wantSession.TranscriptOutput = twriter
	// fmt.Fprintln(twriter, stdOut)

	wantBuf := &bytes.Buffer{}
	gotBuf := &bytes.Buffer{}
	wantBuf.WriteString("hello world\n")
	session := shellspy.NewSession()

	session.Input = strings.NewReader("echo hello world")
	session.Run()
	fmt.Fprint(gotBuf, session.TranscriptOutput)

	want := wantBuf.String()
	got := gotBuf.String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
