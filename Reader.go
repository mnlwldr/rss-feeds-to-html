package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

type Item struct {
	Title  string
	Link   string
	Source string
}

func main() {
	ordered := make(map[string][]Item)

	file, err := os.Open("urls")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		feedUrl := scanner.Text()

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(feedUrl)
		if err != nil {
			log.Printf("Error %s in %s", err, feedUrl)
			continue
		}
		for _, item := range feed.Items {
			if item.PublishedParsed == nil {
				continue
			}
			pubDate := item.PublishedParsed.Format("2006-01-02")
			u, err := url.Parse(feedUrl)
			if err != nil {
				log.Fatal(err)
			}

			p := bluemonday.StripTagsPolicy()
			title := p.Sanitize(item.Title)
			item := Item{Title: title, Link: item.Link, Source: u.Host}
			ordered[pubDate] = append(ordered[pubDate], item)
		}
	}
	keys := make([]string, 0, len(ordered))
	for k := range ordered {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	header()
	for _, d := range keys {
		fmt.Printf("<h2>%s</h2>", d)
		for _, item := range ordered[d] {
			fmt.Printf("<div class=\"item\"><div class=\"link\"><a href=\"%s\" rel=\"noopener nofollow noreferrer\" target=\"_blank\" title=\"%s\">%s</a></div><div class=\"source\">%s</div></div>", item.Link, item.Title, item.Title, item.Source)
		}
	}
}

func header() {
	const HEADER = `<!doctype html><html lang="en"><head><title>RSS Feeds</title><meta name="viewport" content="width=device-width" />
	<style>a, a:visited { color: #000;padding-bottom: .05rem; text-decoration: underline;}body { font: 1em Go, -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif,'Apple Color Emoji', 'Segoe UI Emoji'; max-height: 100%; line-height: 1.4; color: #111; background: #e0dac7; margin: auto; max-width: 960px; }.item {margin-bottom: 1ch;display: grid;grid-template-columns: minmax(max-content, 24%) 1fr;}
  	h1 { font-size: 2em; margin:0;}h2 { margin:1.5em 0;}hr { border:none; border-top: 1px dashed #868277; margin: 40px 0;}.link { grid-column: 2;   grid-row: 1 }.source { grid-column: 1; grid-row: 1 }.source { text-align: left; padding-right: 12px;align-self: center;}.source span {background-color: #fff;font-size: 0.8em;text-transform: uppercase;}.link {white-space: nowrap;overflow: hidden;text-overflow: ellipsis;}
	</style></head><body><main><div class="items">`
	fmt.Print(HEADER)
}

func footer() {
	const FOOTER = `</div></main></body></html>`
	fmt.Print(FOOTER)
}
