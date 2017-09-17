// Tealang runtime REPL tool.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tealang/tea-go/tea/repl"
)

const (
	welcomeText = `Tealang v0.1-alpha
Copyright 2017 Lennart Espe. All rights reserved.`
	replSymbol = ">>> "
)

func main() {
	fmt.Println(welcomeText)

	reader := bufio.NewReader(os.Stdin)
	ui := repl.New()

	for ui.Active {
		fmt.Print(replSymbol)
		input, err := reader.ReadString('\n')
		if err != nil {
			ui.Stop()
		} else {
			output, err := ui.Interpret(strings.TrimRight(input, "\n"))
			if err != nil {
				fmt.Printf("Failed to execute: %v\n", err)
			} else {
				fmt.Println(output)
			}
		}
	}
	fmt.Println()
}
