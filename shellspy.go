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
	"syscall"
)

type Session struct {
	Input      io.Reader
	Output     io.Writer
	Terminal   io.Writer
	Transcript io.Writer
	Port       int
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

func CreateTranscriptFile() (*os.File, error) {

	file, err := os.OpenFile("shellspy.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func RunRemotely(s *Session) error {

	fmt.Fprintf(s.Output, "shellspy is running remotely on port %d and the output file is shellspy.txt\n", s.Port)
	address := fmt.Sprintf("localhost:%d", s.Port)

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT)

	exit_chan := make(chan int)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		for {
			sc := <-signalChanel
			switch sc {
			case syscall.SIGINT:
				conn, _ := listener.Accept()
				go handleConn(conn, s, true)
				exit_chan <- 0
			default:
				conn, _ := listener.Accept()
				go handleConn(conn, s, false)
				exit_chan <- 0
			}
		}
	}()
	exitCode := <-exit_chan
	os.Exit(exitCode)
	return nil
}

func handleConn(c net.Conn, s *Session, terminate bool) {

	if !terminate {
		fmt.Fprintf(c, "welcome to shellspy, output file is shellspy.txt"+"\n")
		s.Input = io.Reader(c)
		s.Start()
	} else {
		fmt.Fprintf(c, "Server has terminated this process"+"\n")
		s.Input = io.Reader(c)
		os.Exit(1)
		s.Start()
	}
}

func (s *Session) Start() {

	scanner := bufio.NewScanner(s.Input)
	buf := &bytes.Buffer{}

	for scanner.Scan() {
		cmd := CommandFromString(scanner.Text())
		input := scanner.Text() + "\n"
		cmd.Stdout = s.Output
		cmd.Stderr = s.Output
		s.Transcript.Write([]byte(input))
		if scanner.Text() == "exit" {
			os.Exit(0)
		}
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			fmt.Fprint(buf, err)
			bbytes := buf.Bytes()
			s.Transcript.Write(bbytes)
		}
	}
}

func CommandFromString(line string) *exec.Cmd {

	trim := strings.TrimSuffix(line, "\n")
	name := strings.Fields(trim)
	args := name[1:]
	cmd := exec.Command(name[0], args...)
	return cmd
}
