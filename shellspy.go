package shellspy

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

type Session struct {
	Input      io.Reader
	Output     io.Writer
	Terminal   io.Writer
	Transcript io.Writer
	Port       int
}

func NewSession(output io.Writer) (*Session, error) {

	s := &Session{}

	file, err := CreateTranscriptFile()
	if err != nil {
		return s, err
	}
	s.Transcript = file
	s.Input = os.Stdin
	s.Terminal = output
	s.Output = io.MultiWriter(s.Terminal, s.Transcript)
	return s, nil
}

func RunCLI(cliArgs []string, output io.Writer) {

	s, err := NewSession(output)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	if len(cliArgs) == 0 {
		fmt.Fprint(s.Output, "shellspy is running locally\n")
		s.Start()
	}

	if len(cliArgs) == 2 {
		fs := flag.NewFlagSet("cmd", flag.ExitOnError)
		portFlag := fs.Int("port", 2000, "-port 3000")

		fs.Parse(cliArgs)

		if portFlag != nil {
			s.Port = *portFlag
			RunRemotely(s)
		}
	}

}

func RunRemotely(s *Session) error {

	fmt.Fprintf(s.Output, "shellspy is running remotely on port %d\n", s.Port)
	address := fmt.Sprintf("localhost:%d", s.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, s)
	}
}

func handleConn(c net.Conn, s *Session) {

	fmt.Fprintf(c, "hello, welcome to shellspy"+"\n")
	s.Input = io.Reader(c)
	s.Start()
}

func (s *Session) Start() {

	scanner := bufio.NewScanner(s.Input)
	for scanner.Scan() {
		cmd := CommandFromString(scanner.Text())
		cmd.Stdout = s.Output
		cmd.Stderr = s.Output
		cmd.Run()
	}
}

func CommandFromString(line string) *exec.Cmd {

	trim := strings.TrimSuffix(line, "\n")
	name := strings.Fields(trim)
	args := name[1:]
	cmd := exec.Command(name[0], args...)
	return cmd
}

func CreateTranscriptFile() (*os.File, error) {

	file, err := os.OpenFile("shellspy.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
