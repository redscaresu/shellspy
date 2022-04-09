package shellspy

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Session struct {
	Input      io.Reader
	Output     io.Writer
	Terminal   io.Writer
	Transcript io.Writer
	Port       string
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

	if len(cliArgs) == 1 {
		fmt.Fprint(s.Output, "shellspy is running locally\n")
		s.Start()
	}
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
