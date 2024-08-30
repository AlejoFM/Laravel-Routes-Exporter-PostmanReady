package postman

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func extractRulesFromRequest(requestPath string) map[string]string {
	content, err := os.ReadFile(requestPath)
	if err != nil {
		fmt.Println("Error reading FormRequest class:", err)
		return nil
	}

	rules := make(map[string]string)
	rulesRegex := regexp.MustCompile(`'(\w+)'\s*=>\s*'(.+?)'`)
	matches := rulesRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		rules[match[1]] = match[2]
	}

	return rules
}

func extractRulesFromValidate(validateContent string) map[string]string {
	rules := make(map[string]string)
	rulesRegex := regexp.MustCompile(`'(\w+)'\s*=>\s*'(.+?)'`)
	matches := rulesRegex.FindAllStringSubmatch(validateContent, -1)
	for _, match := range matches {
		rules[match[1]] = match[2]
	}

	return rules
}

func getDefaultValueForRule(rule string) string {
	if strings.Contains(rule, "required") {
		return "example_value"
	}
	if strings.Contains(rule, "email") {
		return "example@example.com"
	}
	if strings.Contains(rule, "integer") {
		return "0"
	}
	return ""
}
