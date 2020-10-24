package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var hadParseError bool = false
var hadRuntimeError bool = false

func main() {
	arglen := len(os.Args)
	if arglen > 2 {
		fmt.Println("usage: bitracer [script]")
		os.Exit(64)
	} else if arglen == 2 {
		fmt.Println(os.Args)
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func readFileContent(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return string(file)
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome, type '.exit' to exit")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == ".exit\n" {
			fmt.Println("goodbye")
			os.Exit(0)
		}
		tokenizer := newTokenizer(text)

		fmt.Println(tokenizer.scanTokens())
	}
}

func runFile(filename string) {
	content := readFileContent(filename)
	run(content)
}

func run(fileContents string) {
	tokens := newTokenizer(fileContents).scanTokens()
	parser := newParser(tokens)
	expr := parser.parse()
	fmt.Printf("%v\n", expr)

	if hadParseError {
		return
	}

	if hadRuntimeError {
		return
	}

	interpreter := &interpreter{}
	interpreter.interpret(expr)
}
