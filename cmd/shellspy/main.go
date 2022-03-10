package main

import (
	"os"

	"github.com/redscaresu/shellspy"
)

func main() {

	cliArgs := os.Args
	shellspy.RunCLI(cliArgs, os.Stdout)
}
