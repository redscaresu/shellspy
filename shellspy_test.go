package shellspy_test

import (
	"fmt"
	"io/ioutil"
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
	got := shellspy.RunFromCmd(cmd)

	if want != got {
		t.Fatal("something gone wrong")
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

	testString := `echo "hello world"
hello world
ls testdata/
shellspy.txt
`

	writeFile := shellspy.WriteTranscript(testString)

	p, err := ioutil.ReadFile(writeFile.Name())
	if err != nil {
		t.Fatal("something has gone wrong!")
	}
	got := string(p)

	if want != got {
		fmt.Println(want)
		fmt.Println(got)
		t.Fatal("something has gone wrong!")
	}
}
