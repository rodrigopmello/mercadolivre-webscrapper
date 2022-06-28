package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/gocolly/colly/v2"
	"github.com/montanaflynn/stats"
)

// this application can work as one of the agents' actuators
// based on the current intention, an agent can retrieve information about some topic
// sensors - beliefs (topic) - retrieve info from the web - sensor -> neural network - learn and improve an existent plan
// learn and improve an agent strategy, for instance

type Item struct {
	Title       string
	Description string
	Price       float64
	User        string
	// UserRating  string
	Amount int
	New    bool
}

const threshold float64 = 1.0

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

	itens := []Item{}

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	detailCollector := c.Clone()
	// this function main goal is to iterate over
	// the returned list and enter in the links
	// to retrieved the desired information
	c.OnHTML(`a.ui-search-item__group__element.ui-search-link`, func(e *colly.HTMLElement) {
		log.Println("Starting Scraper")
		productLink := e.Attr("href")
		log.Println("Visiting item: ", productLink)
		detailCollector.Visit(productLink)
	})

	detailCollector.OnHTML("div.ui-pdp-container__row.ui-pdp-component-list.pr-16.pl-16", func(e *colly.HTMLElement) {
		log.Println("Extracting product details")
		title := e.ChildText(".ui-pdp-title")
		price := e.ChildText(".price-tag-fraction")
		user := e.ChildText(".ui-pdp-color--BLUE")
		//amount := e.ChildText(".ui-pdp-color--BLACK.ui-pdp-size--XSMALL.ui-pdp-family--REGULAR.ui-pdp-seller__header__subtitle")
		amount := e.ChildText(".ui-pdp-buybox__quantity__available")
		amount = strings.ReplaceAll(amount, "(", "")
		//ui-pdp-buybox__quantity__available
		new := e.ChildText(".ui-pdp-subtitle")
		item := Item{}
		item.Title = title
		priceF, _ := strconv.ParseFloat(price, 32)

		item.Price = float64(priceF) / float64(100)
		item.User = user

		item.Amount, _ = strconv.Atoi(strings.Split(amount, " ")[0])

		itens = append(itens, item)

		item.New = strings.Contains(new, "Novo")
		log.Println(item)
	})

	detailCollector.OnHTML("h1.ui-pdp-title", func(e *colly.HTMLElement) {
		// el := e.Request.Visit(e.Attr("ol"))
		// log.Println(el)
		log.Println("Visiting product", e)

		product := Item{}
		product.Title = e.Text //ui-pdp-title

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	term := flag.String("term", "item", "term to be used during web-scraping")
	flag.Parse()
	fmt.Println(*term)
	site := "https://lista.mercadolivre.com.br/"
	displayMode := "_DisplayType_LF"
	err := c.Visit(site + *(term) + displayMode)
	if err != nil {
		fmt.Println(err)
	}
	results, _ := json.MarshalIndent(itens, "", " ")

	_ = ioutil.WriteFile("results.json", results, 0644)

	data, err := ioutil.ReadFile("results.json")

	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(data, &itens)
	finalData := removeOutliers(&itens)

	perceptions, _ := json.MarshalIndent(finalData, "", " ")

	_ = ioutil.WriteFile("perceptions.json", perceptions, 0644)

}

func findSimilarities(items []Item) {
	for i := 0; i < len(items)/2; i++ {
		for j := 1; j < len(items)/2; j++ {
			similarity := strutil.Similarity(items[i].Title, items[j].Title, metrics.NewHamming())
			log.Println("Similarity ", items[i].Title, items[j].Title, similarity)
		}
	}
}

func extractPrices(items *[]Item) []float64 {
	data := []float64{}
	for _, v := range *items {
		data = append(data, v.Price)
	}

	return data
}

func findQuartile(items *[]Item) (stats.Outliers, error) {
	prices := extractPrices(items)
	q, err := stats.QuartileOutliers(prices)
	if err != nil {
		return stats.Outliers{}, err
	}
	return q, nil
}

func removeOutliers(items *[]Item) []Item {

	data := []float64{}
	for _, v := range *items {
		data = append(data, v.Price)
	}

	std, _ := stats.StandardDeviation(data)
	mean, _ := stats.Mean(data)
	log.Println(std)
	q, _ := stats.Quartile(data)
	log.Println(q)
	cleanedData := []Item{}
	for _, v := range *items {
		if math.Abs(v.Price) > math.Abs(mean-threshold*std) && math.Abs(v.Price) < math.Abs(mean+threshold*std) {
			cleanedData = append(cleanedData, v)
		}
	}

	log.Println(len(cleanedData))
	for _, v := range cleanedData {
		log.Println(v.Price)

	}
	return cleanedData

}
