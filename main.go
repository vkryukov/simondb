package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func parseURL(u string) (*html.Node, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	return html.Parse(resp.Body)
}

func searchName(name string) (*html.Node, error) {
	return parseURL("http://www.imdb.com/find?s=nm&q=" + url.QueryEscape(name))
}

func resultMatcher(n *html.Node) bool {
	if n != nil && n.DataAtom == atom.Td && scrape.Attr(n, "class") == "result_text" {
		return true
	}
	return false
}

func firstChildLinkNode(n *html.Node) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom == atom.A {
			return c
		}
	}
	return nil
}

func firstChildLink(n *html.Node) string {
	return scrape.Attr(firstChildLinkNode(n), "href")
}

func getFilmographyActor(u string) ([]*html.Node, error) {
	root, err := parseURL("http://www.imdb.com" + u)
	if err != nil {
		return nil, err
	}
	return scrape.FindAll(root, filmMatcher), nil
}

func filmMatcher(n *html.Node) bool {
	return n != nil && n.DataAtom == atom.Div &&
		strings.HasPrefix(scrape.Attr(n, "class"), "filmo-row") &&
		// n.Parent != nil && scrape.Attr(n.Parent, "class") == "filmo-category-section" &&
		n.Parent.PrevSibling != nil && n.Parent.PrevSibling.PrevSibling != nil &&
		scrape.Attr(n.Parent.PrevSibling.PrevSibling, "id") == "filmo-head-actor"
}

func readNumber(prompt string, low, high int) int {
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scanln(&input)
		n, err := strconv.ParseInt(input, 10, 32)
		if err != nil || int(n) < low || int(n) > high {
			fmt.Printf("Wrong input: %s\n", input)
		} else {
			return int(n)
		}
	}
}

func firstN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func stringNode(n *html.Node) string {
	return fmt.Sprintf("node: %v id: %s class: %s text: %s\n",
		n.DataAtom, scrape.Attr(n, "id"), scrape.Attr(n, "class"), firstN(scrape.Text(n), 60))
}

func main() {
	fmt.Print("> Actor name: ")
	var input string
	//fmt.Scanln(&input)
	input = "Bruce"
	root, err := searchName(input)
	if err != nil {
		panic(err)
	}
	results := scrape.FindAll(root, resultMatcher)
	fmt.Println("\nTop 10 matching actors:")
	for i, r := range results {
		if i < 10 {
			fmt.Printf("  %d. %s\n", i+1, scrape.Text(r))
		}
	}

	// n := readNumber("> Select actor: ", 1, 10)
	n := 4
	results, err = getFilmographyActor(firstChildLink(results[n-1]))
	if err != nil {
		panic(err)
	}
	for _, r := range results {
		fmt.Printf("%s\n", scrape.Text(r))
		// fmt.Printf("  node: %s\n  parent:  %s\n  parent prev sibling: %s\n",
		// 	stringNode(r), stringNode(r.Parent), stringNode(r.Parent.PrevSibling.PrevSibling))
	}
}
