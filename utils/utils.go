package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/LysetsDal/webscraper-v2/types"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

// ReadJson Read the Docker daemons responses (format json)
func ReadJson(body io.Reader, decodeTo any) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(decodeTo)
}

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJson Write Json with standard header.
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	_ = WriteJson(w, status, map[string]string{"error": err.Error()})
}

func MakeHttpHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if err != nil {
				return
			}
		}
	}
}

// HTML Utils

// GetHTMLBody Extracts the body element from the html node tree
func GetHTMLBody(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("missing <body> in the node tree")
}

// GetHTMLBodyString Extracts the body element from the html node tree as a string
func GetHTMLBodyString(doc *html.Node) (string, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	if body != nil {
		// Convert body node to string (inner HTML)
		var buf bytes.Buffer
		err := html.Render(&buf, body)
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return "", errors.New("missing <body> in the node tree")
}

// PrettyPrint function to format HTML output with indentation
func PrettyPrint(node *html.Node) string {
	var indentLevel int = 0
	var buf strings.Builder
	indent := strings.Repeat("  ", indentLevel)

	switch node.Type {
	case html.ElementNode:
		buf.WriteString(fmt.Sprintf("%s<%s", indent, node.Data))
		for _, attr := range node.Attr {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, attr.Key, attr.Val))
		}
		buf.WriteString(">\n")
	case html.TextNode:
		text := strings.TrimSpace(node.Data)
		if len(text) > 0 {
			buf.WriteString(fmt.Sprintf("%s%s\n", indent, text))
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		buf.WriteString(PrettyPrint(child))
	}

	if node.Type == html.ElementNode {
		buf.WriteString(fmt.Sprintf("%s</%s>\n", indent, node.Data))
	}

	return buf.String()
}
