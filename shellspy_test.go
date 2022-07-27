package shellspy_test

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/redscaresu/shellspy"
)

func TestCommandFromString(t *testing.T) {
	t.Parallel()

	cmdWant := &exec.Cmd{}
	cmdWant.Args = []string{"/bin/echo", "hello", "world"}
	wantCmd := cmdWant.String()
	want := strings.TrimPrefix(wantCmd, " ")
	got := shellspy.CommandFromString(want).String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestCommandFromStringArgs(t *testing.T) {
	t.Parallel()

	input := "ls"
	want := []string{"ls"}
	got := shellspy.CommandFromString(input).Args
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRunCommand(t *testing.T) {
	t.Parallel()

	gotBuf := &bytes.Buffer{}
	s, err := shellspy.NewSession(gotBuf)
	s.Input = strings.NewReader("echo hello world")
	if err != nil {
		t.Fatal("unable to create file")
	}
	tempDir := t.TempDir()

	now := time.Now()
	filename := tempDir + "shellspy-" + now.Format("2006-01-02-15:04:05") + ".txt"
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal("unable to create file")
	}
	s.Transcript = file

	s.Start()
	// fmt.Fprint(wantBuf, s.Terminal)
	got := gotBuf.String()
	want := "welcome to shellspy\n$ hello world\n$ "

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRunWithoutPortFlagRunInteractively(t *testing.T) {
	t.Parallel()

	var flagArgs []string
	buf := &bytes.Buffer{}

	go shellspy.RunCLI(flagArgs, buf)

	for buf.String() == "" {
		time.Sleep(5 * time.Millisecond)
	}

	got := buf.String()

	want := "welcome to shellspy\n$ "

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPortFlagListensOnPort(t *testing.T) {
	t.Parallel()

	flagArgs := []string{"-port", "6666"}

	go shellspy.RunCLI(flagArgs, io.Discard)

	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	for err != nil {
		t.Log("retrying network connection")
		time.Sleep(10 * time.Millisecond)
		conn, err = net.Dial("tcp", "127.0.0.1:6666")
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		t.Fatal("nothing has been returned by the scanner")
	}

	got := scanner.Text()

	want := "welcome to shellspy"

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
