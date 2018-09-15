package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// BoolAnswer check for (y/n) answer
func BoolAnswer(question string) bool {
	for {
		fmt.Printf("%s (y/n): ", question)
		input := bufio.NewScanner(os.Stdin)
		if input.Scan() {
			value := input.Text()
			if strings.HasPrefix(strings.ToLower(value), "y") {
				return true
			}
			if strings.HasPrefix(strings.ToLower(value), "n") {
				return false
			}
		}
	}
}

// StringAnswer check for string answer
func StringAnswer(question, option string) string {
	fmt.Printf("%s (%s): ", question, option)
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		if value := input.Text(); value != "" {
			return strings.ToLower(value)
		}
	}

	return option
}

// OptionAnswer check for string answer
func OptionAnswer(question string, options ...string) string {
	for {
		fmt.Printf("%s (%s): ", question, strings.Join(options, ","))
		input := bufio.NewScanner(os.Stdin)
		if input.Scan() {
			value := strings.ToLower(input.Text())
			for _, option := range options {
				if value == strings.ToLower(option) {
					return value
				}
			}
		}
	}
}
