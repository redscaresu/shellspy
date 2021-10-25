package main

import (
	"bufio"
	"os"
	"shellspy"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	shellspy.CommandFromString(text)
}
