package main

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"os"
	"time"
)

type Product struct {
	Name     string
	Link     string
	Image    string
	Price    string
	Comments []Comment
}

type Comment struct {
	Name    string
	Rating  int
	Content string
}

func main() {
	resp, err := soup.Get("https://www.amazon.in/TVs/b/ref=nav_shopall_sbc_tvelec_television?ie=UTF8&node=1389396031")

	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)
	results := doc.Find("div", "id", "mainResults").FindAll("li", "class", "s-result-item")
	now := time.Now().UTC()

	for _, result := range results {
		link := result.Find("a", "class", "s-access-detail-page").Attrs()["href"]
		name := result.Find("h2", "class", "s-access-title").Text()
		image := result.Find("img", "class", "s-access-image").Attrs()["src"]
		price := result.Find("span", "class", "s-price").Text()

		product := Product{Name: name, Link: link, Image: image, Price: price}

		json.NewEncoder(os.Stdout).Encode(product)
	}

	fmt.Println("Elapsed time: ", time.Since(now))
}
