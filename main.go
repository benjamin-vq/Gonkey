package main

import (
	"fmt"
	"github.com/benja-vq/gonkey/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome %s to the Monkey programming language!\n", usr.Username)
	fmt.Println("Type any commands to begin")
	repl.Start(os.Stdin, os.Stdout)
}
