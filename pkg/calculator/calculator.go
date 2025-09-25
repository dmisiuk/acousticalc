package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Calculation represents a mathematical expression and its result
type Calculation struct {
	Expression string
	Result     float64
}

// Calculator represents a calculator instance
type Calculator struct {
	// Future: could store state, history, etc.
}

// NewCalculator creates a new calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Evaluate evaluates a mathematical expression
func (c *Calculator) Evaluate(expression string) (float64, error) {
	return Evaluate(expression)
}

// Evaluate takes a mathematical expression string and returns the result
func Evaluate(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, errors.New("empty expression")
	}

	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	if len(tokens) == 0 {
		return 0, errors.New("invalid expression")
	}

	result, err := parseAndEvaluate(tokens)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// tokenize converts an expression string into a slice of tokens
func tokenize(expression string) ([]string, error) {
	var tokens []string
	var currentToken strings.Builder

	// Keep track of whether the previous token was an operator or opening parenthesis
	// This helps us identify negative numbers
	previousTokenIsOperator := true

	for _, char := range expression {
		if unicode.IsSpace(char) {
			// If we have a current token, add it to tokens
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			continue
		}

		// Handle operators and parentheses
		if isOperator(char) || char == '(' || char == ')' {
			// If we have a current token, add it to tokens
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}

			// Special handling for minus sign (could be negative number)
			if char == '-' && previousTokenIsOperator {
				// This is likely a negative number, don't add the minus sign yet
				currentToken.WriteRune(char)
			} else {
				// Add the operator or parenthesis as a separate token
				tokens = append(tokens, string(char))
				previousTokenIsOperator = (char == '(' || isOperator(char))
			}
		} else if unicode.IsDigit(char) || char == '.' {
			currentToken.WriteRune(char)
			previousTokenIsOperator = false
		} else {
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	// Add the last token if it exists
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens, nil
}

// isOperator checks if a character is a mathematical operator
func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

// parseAndEvaluate parses and evaluates tokens using the Shunting Yard algorithm approach
func parseAndEvaluate(tokens []string) (float64, error) {
	// For simplicity, we'll implement a recursive descent parser for basic arithmetic
	// with operator precedence

	// Convert to reverse polish notation and evaluate
	values := []float64{}
	operators := []string{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// If token is a number, push it to stack for numbers
		if isNumber(token) {
			val, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token)
			}
			values = append(values, val)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			// Process until we find a matching opening parenthesis
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if len(values) < 2 {
					return 0, errors.New("invalid expression")
				}

				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				result, err := applyOperator(val1, val2, op)
				if err != nil {
					return 0, err
				}

				values = append(values, result)
			}

			// Pop the opening parenthesis
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			} else {
				return 0, errors.New("mismatched parentheses")
			}
		} else if isOperatorString(token) {
			// Process operators according to precedence
			for len(operators) > 0 && operators[len(operators)-1] != "(" &&
				hasPrecedence(operators[len(operators)-1], token) {

				if len(values) < 2 {
					return 0, errors.New("invalid expression")
				}

				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				result, err := applyOperator(val1, val2, op)
				if err != nil {
					return 0, err
				}

				values = append(values, result)
			}
			operators = append(operators, token)
		} else {
			return 0, fmt.Errorf("invalid token: %s", token)
		}
	}

	// Process remaining operators
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" || operators[len(operators)-1] == ")" {
			return 0, errors.New("mismatched parentheses")
		}

		if len(values) < 2 {
			return 0, errors.New("invalid expression")
		}

		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		val2 := values[len(values)-1]
		values = values[:len(values)-1]
		val1 := values[len(values)-1]
		values = values[:len(values)-1]

		result, err := applyOperator(val1, val2, op)
		if err != nil {
			return 0, err
		}

		values = append(values, result)
	}

	if len(values) != 1 {
		return 0, errors.New("invalid expression")
	}

	return values[0], nil
}

// isNumber checks if a string represents a number
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isOperatorString checks if a string is an operator
func isOperatorString(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

// hasPrecedence checks if op1 has higher or equal precedence than op2
func hasPrecedence(op1, op2 string) bool {
	// Parentheses have special handling and should not be compared directly
	if op1 == "(" || op1 == ")" || op2 == "(" || op2 == ")" {
		return false
	}

	// Multiplication and division have higher precedence than addition and subtraction
	if (op1 == "*" || op1 == "/") && (op2 == "+" || op2 == "-") {
		return true
	}

	// If op1 and op2 have the same precedence, we return true for left associativity
	if (op1 == "+" || op1 == "-") && (op2 == "+" || op2 == "-") {
		return true
	}
	if (op1 == "*" || op1 == "/") && (op2 == "*" || op2 == "/") {
		return true
	}

	// Otherwise, op2 has higher precedence
	return false
}

// applyOperator applies an operator to two operands
func applyOperator(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", operator)
	}
}
