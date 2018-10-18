package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var inputScanner = bufio.NewScanner(os.Stdin)

// SetInputScanner allows you to override the default input scanner with your own.
func SetInputScanner(scanner *bufio.Scanner) {
	inputScanner = scanner
}

// GetInput reads input from an input buffer and returns the result as a string.
func GetInput() string {
	inputScanner.Scan()
	return strings.TrimSpace(inputScanner.Text())
}

// Prompt gets input with a prompt and returns the input
func Prompt(prompt string) string {
	fullPrompt := fmt.Sprintf("%s", prompt)
	fmt.Print(fullPrompt + ": ")
	return GetInput()
}
