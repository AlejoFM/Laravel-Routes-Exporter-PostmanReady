package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ExtractFromController(controllerContent []byte, action string, formRequestMap map[string]string) map[string]string {
	functionRegex := regexp.MustCompile(`public function ` + action + `\s*\((.*?)\)\s*{([\s\S]*?)}`)
	functionMatches := functionRegex.FindStringSubmatch(string(controllerContent))

	if len(functionMatches) > 2 {
		functionParams := functionMatches[1]
		functionBody := functionMatches[2]

		formRequestMatch := regexp.MustCompile(`(\w+Request)`).FindStringSubmatch(functionParams)
		if len(formRequestMatch) > 0 {
			formRequestClass := formRequestMatch[1]
			if formRequestNamespace, exists := formRequestMap[formRequestClass]; exists {
				formRequestPath := filepath.Join("Here goes your Laravel project directory", strings.ReplaceAll(formRequestNamespace, "\\", "/")+".php")
				return extractRulesFromRequest(formRequestPath)
			}
		} else {
			validateMatches := regexp.MustCompile(`\$request->validate\(\[(.*?)\]\);`).FindStringSubmatch(functionBody)
			if len(validateMatches) > 1 {
				return extractRulesFromValidate(validateMatches[1])
			}
		}
	}
	return nil
}

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
