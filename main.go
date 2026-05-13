package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "[ERROR]: Missing source code :(")
		os.Exit(1)
	}

	bytes, err := os.ReadFile(os.Args[1])

	mode := ""
	if len(os.Args) > 2 {
		mode = os.Args[2]
	}

	if err != nil {
		println("error leyendo el file: ", err)
	}

	//SCANNEO
	scanner := createScanner()
	tokens, err := scanner.scan(bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error de scanneo:", err)
		os.Exit(1)
	}

	if mode == "--scanner" {
		for _, token := range tokens {
			fmt.Printf("linea %d: %s <%s>\n", token.line, token.name, token.value)
		}
		println()
		os.Exit(0)
	}

	// PARSEO
	parser := createParser(tokens)
	statements, err := parser.parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error de parseo:", err)
		os.Exit(1)
	}

	if mode == "--parser" {
		for _, stmt := range statements {
			fmt.Println(printStmt(stmt, 0))
		}
		os.Exit(0)
	}

	// INTERPRETACIÓN
	interpreter := createInterpreter(statements)
	if err := interpreter.interpret(); err != nil {
		fmt.Fprintln(os.Stderr, "error de ejecución:", err)
		os.Exit(1)
	}

	os.Exit(0)
}
