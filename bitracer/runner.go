package main

import (
	"io/ioutil"
	"log"
	"strings"
)

var hadParseError bool = false
var hadRuntimeError bool = false

func tokenize(fileContents string) []token {
	return newTokenizer(fileContents).scanTokens()
}

func parseStatements(tokens []token) ([]stmt, *parseError) {
	return newParser(tokens).parse()
}

func interpretParsedStatements(stmts []stmt) {
	newInterpreter().interpret(stmts)
}

func readFileContent(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return string(file)
}

// For running files
func staticFileRunner(filename string) {
	if !strings.HasSuffix(filename, ".ob") && !strings.HasSuffix(filename, ".obsidian") {
		panic("file must end in '.ob' or '.obsidian'")
	}

	fileContents := readFileContent(filename)
	tokens := tokenize(fileContents)
	// We can safely throw away this error since it gets reported
	stmts, _ := parseStatements(tokens)

	if hadParseError {
		return
	}

	if hadRuntimeError {
		return
	}

	interpretParsedStatements(stmts)
}

// TODO(HEY) - Left off thinking about to cascade errors up to the repl so that way they can be handled gracefully and not
// always stop execution from happening for statements after.
func replFileRunner(lines string) {
	tokens := tokenize(lines)
	// We can safely throw away this error since it gets reported
	stmts, _ := parseStatements(tokens)

	if hadParseError {
		hadParseError = false
	}

	if hadRuntimeError {
		hadRuntimeError = false
	}

	interpretParsedStatements(stmts)
}
