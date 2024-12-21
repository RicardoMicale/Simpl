package repl

import (
	"bufio"
	"fmt"
	"io"
	"language/lexer"
	"language/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	//	creates a scanner object
	scanner := bufio.NewScanner(in)

	for {
		//	prints the prompt
		fmt.Println(PROMPT)

		//	calls the function to scan the user input
		scanned := scanner.Scan()

		//	if the value entered is null make an early return
		if !scanned {
			return
		}

		//	reads the text from the user
		line := scanner.Text()

		//	creates a Lexer object with the user text as the input
		l := lexer.New(line)

		//	initializes the lexer tokens and loops until the token is an EOF type
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
