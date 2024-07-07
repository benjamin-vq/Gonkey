package main

import (
	"flag"
	"fmt"
	"github.com/benja-vq/gonkey/config"
	"github.com/benja-vq/gonkey/repl"
	"os"
	"os/user"
)

func main() {
	config.Debug = flag.Bool("debug", false, "Prints debugging information during interpreter execution")
	flag.Parse()

	usr, err := user.Current()

	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome %s to the Monkey programming language!\n", usr.Username)
	fmt.Println("Type any commands to begin")
	repl.Start(os.Stdin, os.Stdout)
}
