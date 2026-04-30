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
	if err != nil {
		println("error leyendo el file: ", err)
	}
	scanner := createScanner()
	tokens, err := scanner.scan(bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error de scanneo:", err)
		os.Exit(1)
	}

	for _, token := range tokens {
		fmt.Printf("%s<%s> ", token.name, token.value)
	}
	println()
}
