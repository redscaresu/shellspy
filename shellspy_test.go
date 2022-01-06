package shellspy_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"shellspy"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInput(t *testing.T) {

	input := strings.NewReader("echo hello world")
	want := "hello world\n"

	// var writer *bytes.Buffer
	// foo := bytes.Buffer{}
	// writer = &foo

	// writer := &bytes.Buffer{}
	writer := ""

	session := shellspy.NewSession() //returns a struct
	session.Input = input
	session.Output = writer
	session.Run()

	got := session.Output

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

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
	got := shellspy.RunFromCmd(cmd)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWriteShellScript(t *testing.T) {

	if _, err := os.Stat("shellspy.txt"); err == nil {
		os.Remove("shellspy.txt")
	}

	file, err := os.Open("testdata/shellspy.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	s, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal("something has gone wrong!")
	}
	want := string(s)

	writer := &bytes.Buffer{}
	io.WriteString(writer, `echo "hello world"
hello world
ls testdata/
shellspy.txt
`)

	// 	testString := `echo "hello world"
	// hello world
	// ls testdata/
	// shellspy.txt
	// `

	writeFile := shellspy.WriteTranscript(writer.String())

	p, err := ioutil.ReadFile(writeFile.Name())
	if err != nil {
		t.Fatal("something has gone wrong!")
	}
	got := string(p)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(writer, got))
	}
	os.Remove("shellspy.txt")
}
