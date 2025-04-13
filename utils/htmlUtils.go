package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	T "github.com/LysetsDal/webscraper-v2/types"
	"golang.org/x/net/html"
)

// HTML Utils

func ExtractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += ExtractText(c)
	}
	return result
}

func GetHTMLDoc(r *http.Response) *html.Node {
	doc, err := html.Parse(r.Body)
	if err != nil {
		fmt.Errorf("error parsing response to html: %s", err)
		return nil
	}
	return doc
}

func ParseOpenList(doc *html.Node) ([]T.Entry, error) {
	var entries []T.Entry
	var currentRow []string

	var aux func(*html.Node)
	aux = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			currentRow = []string{}
			// Collect all <td> text values
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					text := ExtractText(c)
					if strings.TrimSpace(text) != "" {
						currentRow = append(currentRow, strings.TrimSpace(text))
					}
				}
			}

			// Add as Entry if it's a real data row (not header)
			if len(currentRow) == 2 && currentRow[1] != "Ekstern venteliste" {
				entries = append(entries, *T.NewEntry(currentRow[0], currentRow[1]))
			}
		}

		// Recurse
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			aux(c)
		}
	}

	aux(doc)

	return entries, nil
}

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

func GetHTMLElement(doc *html.Node, elementTag string) (*html.Node, error) {
	var element *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == elementTag {
			element = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	if element != nil {
		return element, nil
	}

	return nil, fmt.Errorf("missing <%s> in the node tree", elementTag)
}

func GetNestedHTMLElement(doc *html.Node, outerTag string, innerTag string) (*html.Node, error) {

	// Find the Outer element (outerTag)
	outerElement, err := GetHTMLElement(doc, outerTag)
	if err != nil {
		return nil, fmt.Errorf("missing guard element <%s> in the node tree", outerTag)
	}

	// Search for the inner element (innerTag) within outer element (outerTag)
	innerElement, err := GetHTMLElement(outerElement, innerTag)
	if err != nil {
		return nil, fmt.Errorf("missing nested element <%s> inside <%s>", innerTag, outerTag)
	}

	return innerElement, nil
}
