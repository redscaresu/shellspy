package main

import (
	"bufio"
	"fmt"
	"os"
	"shellspy"
	"strings"
)

func main() {

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		cmd, _ := shellspy.CommandFromString(input)
		if strings.HasPrefix(input, "exit") {
			os.Exit(0)
		}

		runCmd := shellspy.RunFromCmd(cmd)
		fmt.Println(runCmd)
	}
}
