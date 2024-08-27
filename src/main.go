package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	rootDir := "Here goes your Laravel proyect directory"

	var items []Item

	// read api.php file
	apiFilePath := filepath.Join(rootDir, "routes", "api.php")
	content, err := os.ReadFile(apiFilePath)
	if err != nil {
		fmt.Println("Error leyendo api.php:", err)
		return
	}

	namespaceMap := make(map[string]string)
	useRegex := regexp.MustCompile(`use\s+([^;]+);`)
	useMatches := useRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range useMatches {
		parts := strings.Split(match[1], "\\")
		controller := parts[len(parts)-1]
		namespaceMap[controller] = match[1]
	}

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
				fmt.Println("Error leyendo controlador:", err)
				continue
			}

			formRequestMap := make(map[string]string)
			useMatches := useRegex.FindAllStringSubmatch(string(controllerContent), -1)
			for _, match := range useMatches {
				if !strings.Contains(match[1], "Illuminate") {
					formRequestClass := filepath.Base(match[1])
					formRequestMap[formRequestClass] = match[1]
				}
			}

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
						rules := extractRulesFromRequest(formRequestPath)
						addPostmanItem(&items, httpMethod, endpoint, rules)
					}
				} else {
					validateMatches := regexp.MustCompile(`\$request->validate\(\[(.*?)\]\);`).FindStringSubmatch(functionBody)
					if len(validateMatches) > 1 {
						rules := extractRulesFromValidate(validateMatches[1])
						addPostmanItem(&items, httpMethod, endpoint, rules)
					} else {
						addPostmanItem(&items, httpMethod, endpoint, nil)
					}
				}
			}
		}
	}

	collection := PostmanCollection{
		Info: Info{
			PostmanID: "uuid.New().String()",
			Name:      "Laravel API Requests",
			Schema:    "https://schema.getpostman.com/json/collection/v2.0.0/collection.json",
		},
		Items: items,
	}

	jsonData, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		fmt.Println("Error al formatear JSON:", err)
		return
	}

	outputFile := "postman_collection.json"
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error al escribir el archivo JSON:", err)
		return
	}

	fmt.Printf("Archivo JSON generado: %s\n", outputFile)
}

func extractRulesFromRequest(requestPath string) map[string]string {
	content, err := os.ReadFile(requestPath)
	if err != nil {
		fmt.Println("Error leyendo FormRequest:", err)
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

func addPostmanItem(items *[]Item, method, endpoint string, rules map[string]string) {
	var rawBody string
	if rules != nil {
		rawBody = "{\n"
		for key := range rules {
			rawBody += fmt.Sprintf("    \"%s\": \"\",\n", key)
		}
		rawBody = strings.TrimSuffix(rawBody, ",\n") + "\n}"
	}

	item := Item{
		Name: endpoint,
		Request: Request{
			Method: method,
			Header: []string{},
			Body: Body{
				Mode: "raw",
				Raw:  rawBody,
				Options: Options{
					Raw: RawOptions{
						Language: "json",
					},
				},
			},
			URL: "{{host}}" + endpoint,
		},
	}
	*items = append(*items, item)
}

type PostmanCollection struct {
	Info  Info   `json:"info"`
	Items []Item `json:"item"`
}

type Info struct {
	PostmanID string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type Item struct {
	Name    string  `json:"name"`
	Request Request `json:"request,omitempty"`
}

type Request struct {
	Method string   `json:"method"`
	Header []string `json:"header"`
	Body   Body     `json:"body,omitempty"`
	URL    string   `json:"url"`
}

type Body struct {
	Mode    string  `json:"mode"`
	Raw     string  `json:"raw"`
	Options Options `json:"options"`
}

type Options struct {
	Raw RawOptions `json:"raw"`
}

type RawOptions struct {
	Language string `json:"language"`
}
