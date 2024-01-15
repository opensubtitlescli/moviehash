package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/opensubtitlescli/moviehash"
)

type cmd struct {
	flag.FlagSet
	h *bool
	v *bool
}

func main() {
	cmd := new(os.Stdout)
	code := cmd.run(os.Args)
	os.Exit(code)
}

func new(output io.Writer) *cmd {
	c := &cmd{}
	c.Init("moviehash", flag.ContinueOnError)
	c.SetOutput(output)
	c.FlagSet.Usage = c.usage
	c.h = c.Bool("h", false, "print help message")
	c.v = c.Bool("v", false, "print moviehash version")
	return c
}

func (c *cmd) run(args []string) int {
	err := c.Parse(args[1:])
	if err != nil {
		return 2
	}

	switch {
	case *c.h:
		c.help()
		return 0
	case *c.v:
		c.version()
		return 0
	case c.NArg() == 0:
		c.help()
		return 2
	default:
		// continue
	}

	for _, p := range c.Args() {
		f, err := os.Open(p)
		if err != nil {
			c.error(err)
			return 2
		}
		defer f.Close()

		s, err := moviehash.Sum(f)
		if err != nil {
			c.error(err)
			return 2
		}

		fmt.Fprintln(c.Output(), s)
	}

	return 0
}

func (c *cmd) help() {
	fmt.Fprintln(c.Output(), "moviehash: calculate the moviehash of the file.")
	c.usage()
}

func (c *cmd) error(err error) {
	fmt.Fprintln(c.Output(), err)
	c.usage()
}

func (c *cmd) usage() {
	fmt.Fprintln(c.Output(), "usage: moviehash [-hv] <path...>")
}

func (c *cmd) version() {
	fmt.Fprintf(c.Output(), "%s\n", moviehash.Version)
}
