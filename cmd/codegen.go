package main

import (
	"fmt"
	"os"

	"github.com/kmio11/codegen/cmd/mock"
)

type Command interface {
	Name() string
	Description() string
	Usage(cmd string)
	Parse(args []string) error
	Execute() int
}

var (
	commands = []Command{
		mock.New(),
	}
)

func usage() {
	indent := "  "
	fmt.Printf(`Usage:

%s%s <command> [arguments]

The commands are:

`, indent, os.Args[0])

	for _, c := range commands {
		fmt.Printf("%s%s  %s", indent, c.Name(), c.Description())
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	for _, c := range commands {
		if os.Args[1] == c.Name() {
			err := c.Parse(os.Args[2:])
			if err != nil {
				c.Usage(os.Args[0])
				os.Exit(1)
			}
			exitCode := c.Execute()
			os.Exit(exitCode)
		}
	}
	usage()
}
