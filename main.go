package main

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"os"
	"time"
)

type Product struct {
	Name    string
	Link    string
	Image   string
	Price   string
	Reviews []Review
}

type Review struct {
	Name    string
	Rating  string
	Content string
}

func (review *Review) parseHtml(raw soup.Root) {
	review.Name = raw.Find("span", "class", "a-profile-name").Text()
	review.Rating = raw.Find("span", "class", "a-icon-alt").Text()
	review.Content = raw.Find("div", "class", "a-expander-content").Text()
}

func (product *Product) getReviews() {
	resp, err := soup.Get(product.Link)

	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)

	reviewsContainer := doc.Find("div", "class", "reviews-content")

	if reviewsContainer.Error != nil {
		return
	}

	rawReviews := reviewsContainer.FindAll("div", "class", "review")
	reviews := []Review{}

	for _, rawReview := range rawReviews {
		contentHolder := rawReview.Find("div", "class", "a-expander-content")

		if contentHolder.Error != nil {
			continue
		}

		review := Review{}
		review.parseHtml(rawReview)

		reviews = append(reviews, review)
	}

	product.Reviews = reviews
}

func (product *Product) parseHtml(result soup.Root) {
	product.Link = result.Find("a", "class", "s-access-detail-page").Attrs()["href"]
	product.Name = result.Find("h2", "class", "s-access-title").Text()
	product.Image = result.Find("img", "class", "s-access-image").Attrs()["src"]
	product.Price = result.Find("span", "class", "s-price").Text()

	product.getReviews()
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
		product := Product{}

		product.parseHtml(result)

		json.NewEncoder(os.Stdout).Encode(product)
	}

	fmt.Println("Elapsed time: ", time.Since(now))
}
