package links

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func FindLinks(htmlDoc []byte) ([]string, error) {

	buf := bytes.NewBuffer(htmlDoc)
	doc, err := html.Parse(buf)
	if err != nil {
		fmt.Errorf("Error finding links: %v\n", err)
		return []string{}, err
	}

	return visit(nil, doc), nil

}

func FindLinksFromUrl(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return []string{}, fmt.Errorf("Error getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

//traverse node tree and extracts links
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

// call breadthFirst(crawl,"url")
func FindLinksBreadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...) // all items of list returned by f append to worklist
			}
		}
	}
}

func Crawl(url string) []string {
	fmt.Println(url)
	list, err := extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	//resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML:%v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue //ignore bad tags
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//pre and post are func args?
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)

	}

	if post != nil {
		post(n)
	}
}
