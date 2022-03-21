package shellspy_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/redscaresu/shellspy"

	"github.com/google/go-cmp/cmp"
)

func TestCommandFromString(t *testing.T) {

	cmdWant := &exec.Cmd{}
	cmdWant.Args = []string{"/bin/echo", "hello", "world"}
	wantCmd := cmdWant.String()
	want := strings.TrimPrefix(wantCmd, " ")
	got := shellspy.CommandFromString(want).String()

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

	cmd = exec.Command("pwd", "-x")
	want = "pwd: illegal option -- x\nusage: pwd [-L | -P]\n"

	_, got = shellspy.RunFromCmd(cmd)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestWriteShellScript(t *testing.T) {

	wantBuf := &bytes.Buffer{}
	gotBuf := &bytes.Buffer{}
	wantBuf.WriteString("hello world\n")
	session, _ := shellspy.NewSession(
		shellspy.WithTranscriptOutput(wantBuf),
	)

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

func TestRunWithoutPortFlagRunInteractively(t *testing.T) {

	buf := &bytes.Buffer{}

	flagArgs := []string{"/var/folders/1v/4mmgcg8s51362djr4g9s9sfw0000gn/T/go-build3590226918/b001/exe/main"}

	shellspy.RunCLI(flagArgs, buf)
	got := buf.String()

	want := "shellspy is running locally\n"

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

// func TestPortFlagStartsNetListener(t *testing.T) {

// 	buf := &bytes.Buffer{}

// 	flagArgs := []string{"/var/folders/1v/4mmgcg8s51362djr4g9s9sfw0000gn/T/go-build3590226918/b001/exe/main", "port", "6666"}

// 	shellspy.RunCLI(flagArgs, buf)
// 	got := buf.String()

// 	want := "shellspy is running remotely on port 6666\n"

// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}
// }
