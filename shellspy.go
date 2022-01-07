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
	"time"
)

type session struct {
	Input  io.Reader
	Output string
}

func RunCLI() {

	if len(os.Args) == 1 {
		fmt.Printf("shellspy is running locally\n")

		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			s := NewSession()
			s.Input = strings.NewReader(input.Text())
			s.Run()
		}
	}

	fmt.Printf("shellspy is running remotely on port %v\n", os.Args[1])

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

func NewSession() session {
	return session{}
}

func (s *session) Run() {

	scanner := bufio.NewScanner(s.Input)
	for scanner.Scan() {
		s.Output = RunServer(scanner.Text())
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

func RunFromCmd(cmd *exec.Cmd) (string, string) {
	var outb bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	cmd.Run()

	stdOut := outb.String()
	stdErr := errb.String()

	return stdOut, stdErr
}

func WriteTranscript(stdOut, stdErr string) os.File {

	now := time.Now()
	filename := "shellspy-" + now.Format("2006-01-02-15:04:05") + ".txt"
	file, err := os.OpenFile(filename,
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
