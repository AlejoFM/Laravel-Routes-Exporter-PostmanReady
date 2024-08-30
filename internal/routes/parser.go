package routes

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlejoFM/Laravel-Routes-Exporter-PostmanReady/internal/postman"
	"github.com/AlejoFM/Laravel-Routes-Exporter-PostmanReady/internal/rules"
)

func ParseRoutes(rootDir string) []postman.Item {
	apiFilePath := filepath.Join(rootDir, "routes", "api.php")
	content, err := os.ReadFile(apiFilePath)
	if err != nil {
		fmt.Println("Error reading api.php:", err)
		return nil
	}

	namespaceMap := extractNamespaces(content)

	routeRegex := regexp.MustCompile(`Route::(get|post|put|delete|patch)\s*\(\s*['"]([^'"]+)['"]\s*,\s*\[\s*([^'"]+)::class\s*,\s*['"]([^'"]+)['"]\s*\]\s*\)`)
	matches := routeRegex.FindAllStringSubmatch(string(content), -1)

	var items []postman.Item
	for _, match := range matches {
		httpMethod, endpoint, controller, action := match[1], match[2], match[3], match[4]

		if namespace, ok := namespaceMap[controller]; ok {
			controllerPath := filepath.Join(rootDir, strings.ReplaceAll(namespace, "\\", "/")+".php")
			controllerContent, err := os.ReadFile(controllerPath)
			if err != nil {
				fmt.Println("Error reading controller:", err)
				continue
			}

			formRequestMap := extractFormRequests(controllerContent)

			rules := rules.ExtractFromController(controllerContent, action, formRequestMap)
			postman.AddItem(&items, httpMethod, endpoint, rules)
		}
	}
	return items
}

func extractNamespaces(content []byte) map[string]string {
	namespaceMap := make(map[string]string)
	useRegex := regexp.MustCompile(`use\s+([^;]+);`)
	useMatches := useRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range useMatches {
		parts := strings.Split(match[1], "\\")
		controller := parts[len(parts)-1]
		namespaceMap[controller] = match[1]
	}
	return namespaceMap
}

func extractFormRequests(controllerContent []byte) map[string]string {
	formRequestMap := make(map[string]string)
	useRegex := regexp.MustCompile(`use\s+([^;]+);`)
	useMatches := useRegex.FindAllStringSubmatch(string(controllerContent), -1)
	for _, match := range useMatches {
		if !strings.Contains(match[1], "Illuminate") {
			formRequestClass := filepath.Base(match[1])
			formRequestMap[formRequestClass] = match[1]
		}
	}
	return formRequestMap
}
