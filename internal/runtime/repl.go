package runtime

// type repl struct {
// 	history  []string
// 	commands []string
// }

// // recompile rebuilds the history into a "file" with explicit delimiting
// // when a user runs a REPL session their variables, classes, etc are
// // always declared.
// func (r *repl) recompile(file []string) string {
// 	return strings.Join(r.commands, "")
// }

// func (r *repl) checkKeyword(text string) bool {
// 	switch text {
// 	case ".exit\n":
// 		os.Exit(0)
// 	case ".history\n":
// 		if len(r.commands) == 0 {
// 			fmt.Println("no commands.")
// 			return true
// 		}
// 		fmt.Println(strings.Join(r.commands, ""))
// 		return true
// 	case ".commands\n":
// 		if len(r.commands) == 0 {
// 			fmt.Println("no commands.")
// 			return true
// 		}
// 		fmt.Println(strings.Join(r.commands, ""))
// 		return true
// 	}
// 	return false
// }

// func (r *repl) start() {
// 	reader := bufio.NewReader(os.Stdin)

// 	fmt.Println("Welcome to obsidian, type '.exit' to exit")
// 	for {
// 		fmt.Print("> ")
// 		text, _ := reader.ReadString('\n')

// 		if r.checkKeyword(text) {
// 			continue
// 		}

// 		r.history = append(r.history, text)

// 		// If the statement fails to parse, do not attempt to interpret at all
// 		_, err := parseStatements(tokenize(text))
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}

// 		// If we aren't invoking a print, don't bother running since we verify above
// 		if !strings.HasPrefix(text, "print") {
// 			r.commands = append(r.commands, text)
// 		} else {
// 			contents := append(r.commands, text)
// 			file := r.recompile(contents)
// 			replFileRunner(file)
// 		}
// 	}
// }
