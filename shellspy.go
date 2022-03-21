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
	Output           io.Writer
	TranscriptOutput io.Writer
	File             *os.File
	Port             string
}

type Option func(*session)

func WithOutput(output io.Writer) Option {
	return func(s *session) {
		s.Output = output
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

func RunCLI(cliArgs []string, w io.Writer) {

	s, err := NewSession(
		WithOutput(w),
	)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	if len(cliArgs) == 1 {
		RunLocally(s, w)
	}

	fs := flag.NewFlagSet("cmd", flag.ContinueOnError)
	fs.Parse(os.Args[1:])

	fmt.Println(os.Args)
	switch os.Args[1] {
	case "port":
		args := fs.Args()
		s.Port = args[1]
		fmt.Println("bollox")
		RunRemotely(s, w)
	}

}

func RunLocally(s *session, w io.Writer) {

	buf := &bytes.Buffer{}
	buf.WriteString("shellspy is running locally\n")
	fmt.Fprint(w, buf)
	s.Output = w
	input := io.Reader(os.Stdin)
	Input(input, s)
}

func RunRemotely(s *session, w io.Writer) error {

	buf := &bytes.Buffer{}
	buf.WriteString("shellspy is running remotely " + s.Port + "\n")
	fmt.Fprint(w, buf)
	s.TranscriptOutput = w

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
	input := io.Reader(c)
	exitStatus := Input(input, s)
	if exitStatus == "0" {
		c.Close()
	}
	c.Close()
}

func Input(input io.Reader, s *session) string {

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		s.Input = strings.NewReader(scanner.Text())
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
	s.Output = writer
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
