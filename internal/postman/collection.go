package postman

import (
	"fmt"
	"strings"
)

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

func addPostmanItem(items *[]Item, method, endpoint string, rules map[string]string) {
	var rawBody string
	if rules != nil {
		rawBody = "{\n"
		for key, rule := range rules {
			value := getDefaultValueForRule(rule)
			rawBody += fmt.Sprintf("    \"%s\": \"%s\",\n", key, value)
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
func NewCollection(name string, items []Item) *PostmanCollection {
	return &PostmanCollection{
		Info: Info{
			PostmanID: "uuid.New().String()",
			Name:      name,
			Schema:    "https://schema.getpostman.com/json/collection/v2.0.0/collection.json",
		},
		Items: items,
	}
}

func AddItem(items *[]Item, method, endpoint string, rules map[string]string) {
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
