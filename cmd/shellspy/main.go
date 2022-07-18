package main

import (
	"os"

	"github.com/redscaresu/shellspy"
)

func main() {

	shellspy.RunCLI(os.Args[1:], os.Stdout)
}
