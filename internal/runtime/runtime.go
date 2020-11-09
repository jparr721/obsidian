package runtime

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jparr721/obsidian/internal/interpreter"
	"github.com/jparr721/obsidian/internal/parser"
	"github.com/jparr721/obsidian/internal/statement"
	"github.com/jparr721/obsidian/internal/tokens"
)

// ObcRT is the runtime entrypoint struct used to kick off the interpreter. It tracks errors during
// file execution and reports them and can run in debug mode to unwind the call stack and show where
// things began to break down
type ObcRT struct {
	didError   bool
	errorStack []error
}

func (o *ObcRT) coreDump(metadata interface{}) {
	fileName := fmt.Sprintf("core_dump_%s.log", time.Now().Format(time.RFC3339))
	metadataStr := fmt.Sprintf("%v", metadata)

	f, err := os.Create(fileName)
	_, err = f.WriteString(metadataStr)
	err = f.Close()
	if err != nil {
		fmt.Printf("Last ditch, failed to dump core to file, puking here. Include this with your issue please")
		fmt.Println(metadataStr)
		return
	}
}

func (o *ObcRT) tokenize(fileContents string) []tokens.Token {
	tokens, err := tokens.NewTokenizer(fileContents).ScanTokens()

	if err != nil {
		o.didError = true
		o.errorStack = append(o.errorStack, err)
	}

	return tokens
}

func (o *ObcRT) parse(tokens []tokens.Token) []statement.Statement {
	parsed, err := parser.NewParser(tokens).Parse()

	if err != nil {
		o.didError = true
		o.errorStack = append(o.errorStack, err)
	}

	return parsed
}

func (o *ObcRT) interpret(statements []statement.Statement) {
	err := interpreter.NewInterpreter().Interpret(statements)

	if err != nil {
		o.didError = true
		o.errorStack = append(o.errorStack, err)
	}
}

func (o *ObcRT) readFileContent(filename string) string {
	if !strings.HasSuffix(filename, ".ob") && !strings.HasSuffix(filename, ".obsidian") {
		log.Fatal("file must end in '.ob' or '.obsidian'")
	}

	file, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	return string(file)
}

func (o *ObcRT) _start(filename string) string {
	return o.readFileContent(filename)
}

// Start starts the runtime object and builds the necessary pieces
func (o *ObcRT) Start(file string, repl bool) {
	if !repl {
		file = o._start(file)
	}

	tokens := o.tokenize(file)

	if o.didError {
		for _, err := range o.errorStack {
			fmt.Println(err)
		}
	}

	statements := o.parse(tokens)

	o.interpret(statements)
}
