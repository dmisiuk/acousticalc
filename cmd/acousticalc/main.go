package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dmisiuk/acousticalc/pkg/calculator"
)

func main() {
	// Check if an expression was provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: acousticalc <expression>")
		fmt.Println("Example: acousticalc \"2 + 3 * 4\"")
		os.Exit(1)
	}

	// Join all arguments to handle expressions with spaces
	expression := strings.Join(os.Args[1:], " ")

	// Evaluate the expression
	result, err := calculator.Evaluate(expression)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print the result
	fmt.Printf("Result: %v\n", result)
}
