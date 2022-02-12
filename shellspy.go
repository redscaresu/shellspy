package shellspy

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type session struct {
	Input            io.Reader
	output           io.Writer
	TranscriptOutput io.Writer
	File             *os.File
	Port             int
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

func NewSession(opts ...Option) *session {

	session := &session{}

	for _, o := range opts {
		o(session)
	}

	return session
}

func RunCLI() {

	file, err := CreateFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	output := &bytes.Buffer{}

	s := NewSession(
		WithOutput(output),
	)

	s.File = file

	local := flag.String("mode", "", "set to run locally")
	port := flag.Int("port", 0, "port number")
	flag.Parse()

	if *local == "" && (*port == 0) {
		fmt.Println("Usage: [ --port int | --mode local ]")
		os.Exit(1)
	}

	if *local == "local" {
		RunLocally(s)
	}

	if *local == "" && (*port >= 1 && *port <= 65535) {
		s.Port = *port
		err := RunRemotely(s)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func RunLocally(s *session) {

	fmt.Printf("shellspy is running locally\n")
	input := bufio.NewScanner(os.Stdin)
	Input(input, s)
}

func RunRemotely(s *session) error {
	fmt.Printf("shellspy is running remotely on port %v\n", s.Port)

	address := "localhost:" + strconv.Itoa(s.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, s)
	}
}

func handleConn(c net.Conn, s *session) {

	input := bufio.NewScanner(c)
	Input(input, s)
	c.Close()
}

func Input(input *bufio.Scanner, s *session) {

	for input.Scan() {
		s.Input = strings.NewReader(input.Text())
		s.Run()
	}
}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

func (s *session) Run() {

	writer := &bytes.Buffer{}
	twriter := &bytes.Buffer{}
	iReader := &bytes.Buffer{}
	fmt.Fprint(iReader, s.Input)
	file := s.File
	// stringIreader := iReader.String()
	//this really sucks
	// replacer := strings.NewReplacer("{", "", "}", "", "&", "", "-1", "", "0", "")
	// stringInput := replacer.Replace(stringIreader)

	input := iReader.String()
	input = strings.TrimPrefix(input, "&{")
	input = strings.TrimSuffix(input, " 0 -1}")
	stdOut := RunServer(input, file)
	s.output = writer
	s.TranscriptOutput = twriter
	fmt.Fprint(writer, stdOut)
	fmt.Fprint(twriter, stdOut)
}

func RunServer(line string, file *os.File) string {

	cmd := CommandFromString(line)

	if strings.HasPrefix(line, "exit") {
		os.Exit(0)
	}

	stdOut, stdErr := RunFromCmd(cmd)
	WriteTranscript(stdOut, stdErr, cmd, file)
	return stdOut
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

func CreateFile() (*os.File, error) {
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
