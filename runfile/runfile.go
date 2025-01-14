package runfile

import (
	"bufio"
	"fmt"
	"language/evaluator"
	"language/lexer"
	"language/object"
	"language/parser"
	"os"
)

func ExecuteFile(fileName string) error {
	//	opens the file
	file, err := os.Open(fileName)

	if err != nil {
		return fmt.Errorf("could not open the file: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	env := object.NewEnvironment()

	inputFile := ""

	for scanner.Scan() {
		line := scanner.Text()

		inputFile += line
	}

	l := lexer.New(inputFile)
	p := parser.New(l)

	program := p.ParserProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		// return nil
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}

	return nil
}

func printParserErrors(errors []string) {
	fmt.Println("Errors:")
	for i, msg := range errors {
		fmt.Printf("\t%d: %s\n", i+1, msg)
	}
}
