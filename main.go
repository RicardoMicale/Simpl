package main

import (
	"fmt"
	"language/repl"
	"os"
	"os/user"
)

const LANGUAGE_NAME = "Simpl"

func main() {
	//	gets the current user from the os and any error there may be
	user, err := user.Current()

	//	checks if there is any error
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to %s, %s!\n", LANGUAGE_NAME, user.Name)
	fmt.Println("Write your code below (you can write quit to exit the program)")
	repl.Start(os.Stdin, os.Stdout)

}

// func main() {
// 	user, err := user.Current()

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Welcome to %s, %s\n", LANGUAGE_NAME, user.Name)

// 	if len(os.Args) > 1 {
// 		fileName := os.Args[1]
// 		err := runfile.ExecuteFile(fileName)

// 		if err != nil {
// 			fmt.Printf("Error executing file %s: %s\n", fileName, err)
// 		}

// 	} else {
// 		fmt.Println("No file provided. Starting REPL...")
// 		fmt.Println("Write your code below:")
// 		repl.Start(os.Stdin, os.Stdout)
// 	}
// }
