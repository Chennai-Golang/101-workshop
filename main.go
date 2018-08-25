package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/anaskhan96/soup"
)

// A Product represents a product on amazon
type Product struct {
	Name    string
	Link    string
	Image   string
	Price   string
	Reviews []Review
}

// A Review represents a review on amazon
type Review struct {
	Name    string
	Rating  string
	Content string
}

func (review *Review) parseReviews(raw soup.Root) error {
	contentHolder := raw.Find("div", "class", "a-expander-content")

	if contentHolder.Error != nil {
		return contentHolder.Error
	}

	review.Name = raw.Find("span", "class", "a-profile-name").Text()
	review.Rating = raw.Find("span", "class", "a-icon-alt").Text()
	review.Content = contentHolder.Text()

	return nil
}

func (product *Product) getReviews() {
	now := time.Now().UTC()
	resp, err := soup.Get(product.Link)
	fmt.Println("Fetching time: ", time.Since(now))

	now = time.Now().UTC()

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
		review := Review{}
		err := review.parseReviews(rawReview)

		if err == nil {
			reviews = append(reviews, review)
		}
	}

	product.Reviews = reviews
	fmt.Println("Review parsing time: ", time.Since(now))
}

func parseProducts(result soup.Root, resultChan chan Product) {
	product := Product{}

	product.Link = result.Find("a", "class", "s-access-detail-page").Attrs()["href"]
	product.Name = result.Find("h2", "class", "s-access-title").Text()
	product.Image = result.Find("img", "class", "s-access-image").Attrs()["src"]
	product.Price = result.Find("span", "class", "s-price").Text()

	product.getReviews()

	resultChan <- product
}

func main() {
	runtime.GOMAXPROCS(16)
	now := time.Now().UTC()

	resp, err := soup.Get("https://www.amazon.in/TVs/b/ref=nav_shopall_sbc_tvelec_television?ie=UTF8&node=1389396031")

	fmt.Println("Main fetch time: ", time.Since(now))
	now = time.Now().UTC()

	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)
	results := doc.Find("div", "id", "mainResults").FindAll("li", "class", "s-result-item")

	resultsChan := make(chan Product)
	for _, result := range results {
		go parseProducts(result, resultsChan)
	}

	products := []Product{}
	for range results {
		products = append(products, <-resultsChan)
	}

	json.NewEncoder(os.Stdout).Encode(products)

	fmt.Println("Elapsed time: ", time.Since(now))
}
