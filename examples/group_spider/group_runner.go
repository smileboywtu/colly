package main

import (
	"github.com/asciimoo/colly"
	"fmt"
	"time"
)

/**
	run two runner in a group

 */

// spider a implements runnable
func basic_spider() {

	c := colly.NewCollector()

	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	c.AllowedDomains = []string{"hackerspaces.org", "wiki.hackerspaces.org"}

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://hackerspaces.org/")
}

func simple_spider() {
	url := "https://httpbin.org/delay/2"

	// Instantiate default collector
	c := colly.NewCollector()

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	// Before making a request print "Starting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting", r.URL, time.Now())
	})

	// After making a request print "Finished ..."
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL, time.Now())
	})

	// Start scraping in four threads on https://httpbin.org/delay/2
	for i := 0; i < 4; i++ {
		go c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Start scraping on https://httpbin.org/delay/2
	c.Visit(url)
	// Wait until threads are finished
	c.Wait()
}

func main() {

	g := colly.NewGroup("project1", 10)

	g.AddSpider(basic_spider)
	g.AddSpider(simple_spider)

	g.RunPending()

	g.Wait()

	g.RunForever()

	g.Wait()
}
