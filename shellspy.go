package shellspy

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func RunCli() {
	os.Remove("shellspy.txt")
	listener, err := net.Listen("tcp", "localhost:31359")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func RunServer(input string) {

	cmd, err := CommandFromString(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.HasPrefix(input, "exit") {
		os.Exit(0)
	}

	stdOut, _ := RunFromCmd(cmd)
	WriteTranscript(stdOut)
	fmt.Println(stdOut)
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		RunServer(input.Text())
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
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

	file, err := os.OpenFile("shellspy.txt",
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
