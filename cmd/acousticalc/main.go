package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmisiuk/acousticalc/pkg/calculator"
	"github.com/dmisiuk/acousticalc/pkg/tui"
)

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		// If arguments provided, use CLI mode for backward compatibility
		runCLI()
		return
	}

	// Otherwise, launch TUI mode
	runTUI()
}

// runCLI runs the calculator in command-line mode
func runCLI() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: acousticalc <expression>")
		fmt.Println("Example: acousticalc \"2 + 3 * 4\"")
		fmt.Println("Or run without arguments to start the interactive TUI")
		os.Exit(1)
	}

	// Check for help flags
	arg := os.Args[1]
	if arg == "--help" || arg == "-h" || arg == "help" {
		fmt.Println("AcoustiCalc - A terminal-based calculator")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  acousticalc <expression>    Evaluate expression in CLI mode")
		fmt.Println("  acousticalc                 Start interactive TUI mode")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  acousticalc \"2 + 3 * 4\"")
		fmt.Println("  acousticalc \"(10 + 5) / 3\"")
		fmt.Println()
		fmt.Println("Supported operations: +, -, *, /, (, )")
		fmt.Println("TUI mode supports mouse clicks and keyboard navigation")
		return
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

// runTUI runs the calculator in TUI mode
func runTUI() {
	// Create the TUI model
	model := tui.NewModel()

	// Create the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
