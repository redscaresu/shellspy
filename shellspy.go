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
	Input            io.Reader
	Output           io.Writer
	TranscriptOutput io.Writer
	File             *os.File
}

func RunCLI() {

	file := CreateFile()
	s := NewSession()
	s.File = file

	if len(os.Args) == 1 {
		fmt.Printf("shellspy is running locally\n")

		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			s.Input = strings.NewReader(input.Text())
			s.Run()
		}
	}

	fmt.Printf("shellspy is running remotely on port %v\n", os.Args[1])

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
		go handleConn(conn, s)
	}
}

func NewSession() session {
	return session{}
}

func (s *session) Run() {

	writer := &bytes.Buffer{}
	twriter := &bytes.Buffer{}

	scanner := bufio.NewScanner(s.Input)
	for scanner.Scan() {
		file := s.File
		stdOut := RunServer(scanner.Text(), file)
		s.Output = writer
		s.TranscriptOutput = twriter
		fmt.Fprint(writer, stdOut)
		fmt.Fprint(twriter, stdOut)
	}
}

func RunServer(line string, file *os.File) string {

	cmd, err := CommandFromString(line)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.HasPrefix(line, "exit") {
		os.Exit(0)
	}

	stdOut, stdErr := RunFromCmd(cmd)
	WriteTranscript(stdOut, stdErr, cmd, file)
	return stdOut
}

func handleConn(c net.Conn, s session) {

	input := bufio.NewScanner(c)
	for input.Scan() {
		s.Input = strings.NewReader(input.Text())
		s.Run()
	}
	c.Close()
}

func CommandFromString(line string) (*exec.Cmd, error) {
	trim := strings.TrimSuffix(line, "\n")
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

func CreateFile() *os.File {
	now := time.Now()
	filename := "shellspy-" + now.Format("2006-01-02-15:04:05") + ".txt"
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	return file
}

func WriteTranscript(stdOut, stdErr string, cmd *exec.Cmd, file *os.File) os.File {

	if _, err := file.WriteString(cmd.String()); err != nil {
		log.Println(err)
	}
	file.WriteString("\n")
	if _, err := file.WriteString(stdOut); err != nil {
		log.Println(err)
	}

	return *file
}
