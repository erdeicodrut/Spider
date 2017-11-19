package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

func main() {

	duration := flag.Int("duration", 10, "How long should it crawl for")
	link := flag.String("link", "", "The link you want to crawl")
	flag.Parse()

	links := make(chan string)
	exLinks := make([]string, 1000)

	go func() {
		for {
			select {
			case link := <-links:
				if found(link, exLinks) {
					continue
				}
				go scrape(link, links)
				exLinks = append(exLinks, link)
			}
		}
	}()

	scrape(*link, links)

	time.Sleep(time.Duration(int64(*duration)) * time.Second)

}

func scrape(link string, links chan string) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("Failed to load %v, \nERROR: %v", link, err)
		return
	}

	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" {
					if len(a.Val) < 4 {
						break
					}
					if a.Val[:4] == "http" {
						fmt.Printf("\nFound %v, \nOn %v\n", a.Val, link)
						links <- a.Val
						break
					}
				}
			}
		}
	}
}

func found(who string, where []string) bool {
	for _, link := range where {
		if link == who {
			return true
		}
	}
	return false
}
