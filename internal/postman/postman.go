package postman

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GeneratePostmanItems(rootDir string) []Item {
	var items []Item

	apiFilePath := filepath.Join(rootDir, "routes", "api.php")
	content, err := os.ReadFile(apiFilePath)
	if err != nil {
		fmt.Println("Error reading api.php:", err)
		return nil
	}

	namespaceMap := extractNamespaces(string(content))
	routeRegex := regexp.MustCompile(`Route::(get|post|put|delete|patch)\s*\(\s*['"]([^'"]+)['"]\s*,\s*\[\s*([^'"]+)::class\s*,\s*['"]([^'"]+)['"]\s*\]\s*\)`)

	matches := routeRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		httpMethod := match[1]
		endpoint := match[2]
		controller := match[3]
		action := match[4]

		if namespace, ok := namespaceMap[controller]; ok {
			controllerPath := filepath.Join(rootDir, strings.ReplaceAll(namespace, "\\", "/")+".php")
			controllerContent, err := os.ReadFile(controllerPath)
			if err != nil {
				fmt.Println("Error reading controller:", err)
				continue
			}

			rules := extractRulesFromController(controllerContent, action, rootDir)
			addPostmanItem(&items, httpMethod, endpoint, rules)
		}
	}

	return items
}

func extractNamespaces(content string) map[string]string {
	namespaceMap := make(map[string]string)
	useRegex := regexp.MustCompile(`use\s+([^;]+);`)
	useMatches := useRegex.FindAllStringSubmatch(content, -1)
	for _, match := range useMatches {
		parts := strings.Split(match[1], "\\")
		controller := parts[len(parts)-1]
		namespaceMap[controller] = match[1]
	}
	return namespaceMap
}

func extractRulesFromController(controllerContent []byte, action, rootDir string) map[string]string {
	formRequestMap := extractNamespaces(string(controllerContent))
	functionRegex := regexp.MustCompile(`public function ` + action + `\s*\((.*?)\)\s*{([\s\S]*?)}`)
	functionMatches := functionRegex.FindStringSubmatch(string(controllerContent))

	if len(functionMatches) > 2 {
		functionParams := functionMatches[1]
		functionBody := functionMatches[2]

		formRequestMatch := regexp.MustCompile(`(\w+Request)`).FindStringSubmatch(functionParams)
		if len(formRequestMatch) > 0 {
			formRequestClass := formRequestMatch[1]
			if formRequestNamespace, exists := formRequestMap[formRequestClass]; exists {
				formRequestPath := filepath.Join(rootDir, strings.ReplaceAll(formRequestNamespace, "\\", "/")+".php")
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

func FormatJSON(collection PostmanCollection) ([]byte, error) {
	return json.MarshalIndent(collection, "", "  ")
}
