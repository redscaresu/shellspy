package shellspy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type session struct {
	Input  io.Reader
	Output string
}

func RunCLI() {

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "missing PORT arg\n")
		os.Exit(1)
	}

	os.Remove("shellspy.txt")
	address := "localhost:" + os.Args[1]

	listener, err := net.Listen("tcp", address)
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

func RunServer(input string) string {

	cmd, err := CommandFromString(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.HasPrefix(input, "exit") {
		os.Exit(0)
	}

	stdOut := RunFromCmd(cmd)
	WriteTranscript(stdOut)
	return stdOut
}

func handleConn(c net.Conn) {

	input := bufio.NewScanner(c)
	for input.Scan() {
		RunServer(input.Text())
	}
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

func RunFromCmd(cmd *exec.Cmd) string {
	var outb bytes.Buffer
	cmd.Stdout = &outb

	cmd.Run()

	stdOut := outb.String()
	return stdOut
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

func NewSession() session {
	return session{}
}

func (s *session) Run() {

	scanner := bufio.NewScanner(s.Input)
	for scanner.Scan() {
		s.Output = RunServer(scanner.Text())
	}
}
