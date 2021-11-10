package main

import (
	"bufio"
	"fmt"
	"os"
	"shellspy"
	"strings"
)

func main() {

	os.Remove("transcript.txt")

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		shellspy.WriteTranscript(input)
		cmd, _ := shellspy.CommandFromString(input)
		if strings.HasPrefix(input, "exit") {
			os.Exit(0)
		}

		stdOut, _ := shellspy.RunFromCmd(cmd)
		shellspy.WriteTranscript(stdOut)
		fmt.Println(stdOut)
	}
}
