package shellspy

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type session struct {
	Input            io.Reader
	output           io.Writer
	TranscriptOutput io.Writer
	File             *os.File
	Port             string
}

type Option func(*session)

func WithOutput(output io.Writer) Option {
	return func(s *session) {
		s.output = output
	}
}

func WithTranscriptOutput(TranscriptOutput io.Writer) Option {
	return func(s *session) {
		s.TranscriptOutput = TranscriptOutput
	}
}

func NewSession(opts ...Option) (*session, error) {

	session := &session{}

	for _, o := range opts {
		o(session)
	}

	file, err := CreateTranscriptFile()
	if err != nil {
		return session, err
	}
	session.File = file

	return session, nil
}

func RunCLI() {

	output := os.Stdout

	s, err := NewSession(
		WithOutput(output),
	)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: [ --port int | --mode local ]")
		os.Exit(1)
	}

	cmd := flag.NewFlagSet("cmd", flag.ExitOnError)
	cmd.Bool("help", false, "Print this help message")
	cmd.Bool("h", false, "Print this help message")

	switch os.Args[1] {
	case "help":
		fmt.Println("Usage: [ port int | local | help | h ]")
	case "h":
		fmt.Println("Usage: [ port int | local | help | h ]")
	case "port":
		cmd.Parse(os.Args[1:])
		s.Port = cmd.Arg(1)
		RunRemotely(s)
	case "local":
		cmd.Parse(os.Args[2:])
		RunLocally(s)
	default:
		fmt.Println("Usage: [ port int | local | help | h ]")
		os.Exit(1)
	}
}

func RunLocally(s *session) {

	buf := &bytes.Buffer{}
	fmt.Fprint(buf, "shellspy is running locally\n")
	s.output = buf
	fmt.Print(s.output)
	input := bufio.NewScanner(os.Stdin)
	Input(input, s)
}

func RunRemotely(s *session) error {
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, "shellspy is running locally "+s.Port+"\n")
	s.output = buf
	fmt.Print(s.output)

	address := "localhost:" + s.Port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, s)

		<-killSignal
		fmt.Println("\nconnection terminated by server!")
		listener.Close()
	}
}

func handleConn(c net.Conn, s *session) {

	fmt.Fprintf(c, "hello, welcome to shellspy"+"\n")
	input := bufio.NewScanner(c)
	exitStatus := Input(input, s)
	if exitStatus == "0" {
		c.Close()
	}
	c.Close()
}

func Input(input *bufio.Scanner, s *session) string {

	for input.Scan() {
		s.Input = strings.NewReader(input.Text())
		exitStatus := s.Run()
		if exitStatus == "0" {
			return "0"
		}
	}
	return ""
}

func (s *session) Run() string {

	writer := &bytes.Buffer{}
	twriter := &bytes.Buffer{}
	iReader := &bytes.Buffer{}
	fmt.Fprint(iReader, s.Input)
	file := s.File
	input := iReader.String()
	input = strings.TrimPrefix(input, "&{")
	input = strings.TrimSuffix(input, " 0 -1}")
	stdOut, exitStatus := RunServer(input, file)
	if exitStatus == "0" {
		return "0"
	}
	s.output = writer
	s.TranscriptOutput = twriter
	fmt.Fprint(writer, stdOut)
	fmt.Fprint(twriter, stdOut)
	return ""
}

func RunServer(line string, file *os.File) (string, string) {

	cmd := CommandFromString(line)

	if strings.HasPrefix(line, "exit") {
		return "", "0"
	}

	stdOut, stdErr := RunFromCmd(cmd)
	WriteTranscript(stdOut, stdErr, cmd, file)
	return stdOut, ""
}

func CommandFromString(line string) *exec.Cmd {
	trim := strings.TrimSuffix(line, "\n")
	name := strings.Fields(trim)
	args := name[1:]
	join := strings.Join(args, " ")
	cmd := exec.Command(name[0], join)
	return cmd
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

func CreateTranscriptFile() (*os.File, error) {
	now := time.Now()
	filename := ".shellspy-" + now.Format("2006-01-02-15:04:05") + ".txt"
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteTranscript(stdOut, stdErr string, cmd *exec.Cmd, file *os.File) os.File {

	if _, err := file.WriteString(cmd.String()); err != nil {
		err = fmt.Errorf("unable to write cmd to disk due to error; %w", err)
		file.WriteString(err.Error())
	}

	file.WriteString("\n")

	if stdErr != "" {
		file.WriteString(stdErr)
	}

	if _, err := file.WriteString(stdOut); err != nil {
		err = fmt.Errorf("unable to write stdOut to disk due to error; %w", err)
		file.WriteString(err.Error())
	}

	return *file
}
