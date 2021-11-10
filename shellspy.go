package shellspy

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCli() {
	os.Remove("transcript.txt")

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		WriteTranscript(input)
		cmd, _ := CommandFromString(input)
		if strings.HasPrefix(input, "exit") {
			os.Exit(0)
		}

		stdOut, _ := RunFromCmd(cmd)
		WriteTranscript(stdOut)
		fmt.Println(stdOut)
	}
}

func CommandFromString(input string) (*exec.Cmd, error) {
	trim := strings.TrimSuffix(input, "\n")
	name := strings.Fields(trim)
	args := name[1:]
	join := strings.Join(args, " ")
	cmd := exec.Command(name[0], join)
	return cmd, nil
}

func RunFromCmd(cmd *exec.Cmd) (string, string) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	cmd.Run()

	stdOut := outb.String()
	stdErr := errb.String()

	return stdOut, stdErr
}

func WriteTranscript(stdOut string) os.File {

	file, err := os.OpenFile("transcript.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()
	if _, err := file.WriteString(stdOut); err != nil {
		log.Println(err)
	}

	return *file
}
