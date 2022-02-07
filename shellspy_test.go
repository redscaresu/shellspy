package shellspy_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"shellspy"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCommandFromString(t *testing.T) {

	input := "echo hello world\n"
	want := "/bin/echo hello world"
	got := shellspy.CommandFromString(input).String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
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
	tempDir := t.TempDir()

	now := time.Now()
	filename := tempDir + "shellspy-" + now.Format("2006-01-02-15:04:05") + ".txt"
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal("unable to create file")
	}

	session.File = file
	session.Run()
	fmt.Fprint(gotBuf, session.TranscriptOutput)

	want := wantBuf.String()
	got := gotBuf.String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
