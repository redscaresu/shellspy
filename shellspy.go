package shellspy

import (
	"bufio"
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
}

type Server struct {
	Port int
	C    net.Conn
}

func RunCLI(cliArgs []string, output io.Writer) {

	if len(cliArgs) == 0 {
		RunLocally(output)
	}

	if len(cliArgs) == 2 {
		fs := flag.NewFlagSet("cmd", flag.ExitOnError)
		portFlag := fs.Int("port", 2000, "-port 3000")

		fs.Parse(cliArgs)

		if portFlag != nil {
			RunRemotely(*portFlag)
		}
	}
}

func NewSession(output io.Writer) (*Session, error) {

	s := &Session{}

	file, err := CreateTranscriptFile()
	if err != nil {
		return nil, err
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

func RunRemotely(port int) error {

	fmt.Printf("shellspy is running remotely on port %v and the output file is shellspy.txt\n", port)
	address := fmt.Sprintf("localhost:%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	for {

		killSignal := make(chan os.Signal, 1)
		signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-killSignal
			os.Exit(0)
		}()

		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleConn(conn)

	}
}

func handleConn(c net.Conn) {
	s, err := NewSession(c)
	if err != nil {
		fmt.Fprint(c, err)
		c.Close()
		fmt.Printf("connection is closed due to: %v", err)
	}

	fmt.Printf("new connection established %s and the file is shellspy.txt ", c.RemoteAddr())
	s.Input = c
	s.Start()
}

func RunLocally(output io.Writer) {
	s, err := NewSession(output)
	if err != nil {
		fmt.Println("cannot create transcript file")
		os.Exit(1)
	}
	s.Start()
}

func (s *Session) Start() {

	fmt.Fprintln(s.Output, "welcome to shellspy")
	fmt.Fprintf(s.Output, "$ ")
	scanner := bufio.NewScanner(s.Input)

	for scanner.Scan() {
		cmd := CommandFromString(scanner.Text())
		input := scanner.Text() + "\n"
		cmd.Stdout = s.Output
		cmd.Stderr = s.Output
		fmt.Fprint(s.Transcript, input)
		if scanner.Text() == "exit" {
			os.Exit(0)
		}
		err := cmd.Run()
		if err != nil {
			fmt.Fprint(s.Output, err)
		}
		fmt.Fprintf(s.Output, "$ ")
	}
}

func CommandFromString(line string) *exec.Cmd {

	trim := strings.TrimSuffix(line, "\n")
	name := strings.Fields(trim)
	args := name[1:]
	cmd := exec.Command(name[0], args...)
	return cmd
}
