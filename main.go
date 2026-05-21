package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) < 2 {
		runRepl()
		return
	}

	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error leyendo el file:", err)
		os.Exit(1)
	}

	mode := ""
	if len(os.Args) > 2 {
		mode = os.Args[2]
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

	// RESOLUCIÓN
	resolver := createResolver()
	if err := resolver.resolve(statements); err != nil {
		fmt.Fprintln(os.Stderr, "error de resolución:", err)
		os.Exit(1)
	}

	// INTERPRETACIÓN
	interpreter := createInterpreter(statements, resolver.distances)
	if err := interpreter.interpret(); err != nil {
		fmt.Fprintln(os.Stderr, "error de ejecución:", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func runRepl() {
	fmt.Println("escribí 'exit' para salir")

	sc := createScanner()
	res := createResolver()
	interp := createInterpreter([]Stmt{}, res.distances)

	rl, err := readline.New("> ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error iniciando REPL:", err)
		os.Exit(1)
	}
	defer rl.Close()

	var accumulated string

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println()
			break
		}

		if accumulated == "" && line == "exit" {
			break
		}

		accumulated += line + "\n"

		open := 0
		for _, ch := range accumulated {
			if ch == '{' {
				open++
			} else if ch == '}' {
				open--
			}
		}

		if open > 0 {
			rl.SetPrompt("  ")
			continue
		}
		rl.SetPrompt("> ")

		input := accumulated
		accumulated = ""

		if strings.TrimSpace(input) == "" {
			continue
		}

		tokens, err := sc.scan([]byte(input))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error de scanneo:", err)
			continue
		}

		parser := createParser(tokens)
		statements, err := parser.parse()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error de parseo:", err)
			continue
		}

		if err := res.resolve(statements); err != nil {
			fmt.Fprintln(os.Stderr, "error de resolución:", err)
			continue
		}

		interp.statements = append(interp.statements, statements...)
		for _, stmt := range statements {
			if execErr := interp.execute(stmt); execErr != nil {
				fmt.Fprintln(os.Stderr, "error de ejecución:", execErr)
			}
		}
	}
}
