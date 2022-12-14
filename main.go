package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gocolly/colly"
)

type Flat struct {
	ID        string `json: "id"`
	Title     string `json: "title"`
	Desc      string `json: "desc"`
	Price     string `json: "price"`
	Space     string `json: "space"`
	META      string `json: "meta"`
	PUBLISHED string `json: "published"`
	URL       string `json: "url"`
}

func main() {
	log.Println("Flatty ... 🕵🏾")
	flats := make([]Flat, 0)
	list := ".itemlist"

	// TODO add domains to env
	collector := colly.NewCollector(
		colly.AllowedDomains("ebay-kleinanzeigen.de", "www.ebay-kleinanzeigen.de"),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnHTML(list, func(e *colly.HTMLElement) {
		// class selector -> div.class
		// id selector -> div#id
		// anchor tag -> a[href]

		listItem := ".ad-listitem"
		e.ForEach(listItem, func(_ int, el *colly.HTMLElement) {
			// desc := el.ChildText(".ad-listitem .aditem-main--middle h2")
			// desc := el.ChildAttr(".ad-listitem .aditem-main--middle h2")
			flats = addFlat(el, flats)
		})

	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Error scraping:", err)
	})

	collector.Visit("https://www.ebay-kleinanzeigen.de/s-wohnung-mieten/friedrichshain-kreuzberg/c203l26918r10")

	writeJSON(flats)
}

func addFlat(el *colly.HTMLElement, flats []Flat) []Flat {
	// ATTRIBUTES
	id := el.ChildAttr(".ad-listitem .aditem", "data-adid")
	url := el.ChildAttr(".ad-listitem .aditem-main--middle h2>a", "href")

	// CONTENT
	title := el.ChildText(".ad-listitem .aditem-main--middle h2>a")
	price := el.ChildText(".ad-listitem .aditem-main--middle--price")
	space := el.ChildText(".ad-listitem .aditem-main--bottom .simpletag.tag-small")
	meta := el.ChildText(".ad-listitem .aditem-main--top--left")
	published := el.ChildText(".ad-listitem .aditem-main--top--right")
	desc := el.ChildText(".ad-listitem .aditem-main--middle--description")

	flat := Flat{
		ID:        id,
		Price:     price,
		Desc:      desc,
		Title:     title,
		Space:     space,
		META:      meta,
		PUBLISHED: published,
		URL:       url,
	}

	flats = append(flats, flat)
	return flats
}

func writeJSON(data []Flat) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
	}

	_ = ioutil.WriteFile("flats.json", file, 0644)
	log.Println("Successfuly stored data")
	log.Println("Successfuly stored data")
}
