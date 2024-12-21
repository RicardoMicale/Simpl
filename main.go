package main

import (
	"fmt"
	"language/repl"
	"os"
	"os/user"
)

const LANGUAGE_NAME = "TBD"

func main() {
	//	gets the current user from the os and any error there may be
	user, err := user.Current()

	//	checks if there is any error
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to %s, %s!\n", LANGUAGE_NAME, user.Name)
	fmt.Println("Write your code below")
	repl.Start(os.Stdin, os.Stdout)

}
