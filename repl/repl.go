package repl

import (
	"bufio"
	"fmt"
	"io"
	"language/evaluator"
	"language/lexer"
	"language/object"
	"language/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}
		line := scanner.Text()

		if line == "quit" {
			fmt.Fprintln(out, "Goodbye!")
			break
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParserProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
