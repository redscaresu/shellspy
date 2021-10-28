package main

import (
	"bufio"
	"os"
	"shellspy"
	"strings"
)

func main() {

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		shellspy.CommandFromString(input)

		if strings.HasPrefix(input, "exit") {
			os.Exit(0)
		}
	}
}
