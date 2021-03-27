package main

import (
	"encoding/json"
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
	Title       string //colocar anotacao de len
	Description string
	Price       float64
	User        string
	// UserRating  string
	Amount int
	New    bool
}

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

		log.Println("Starting crawler")
		productLink := e.Attr("href")
		//e.Request.Visit(productLink)
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

	// term := "lg k10"
	// site := "https://lista.mercadolivre.com.br/"
	// displayMode := "_DisplayType_LF"
	// // statisticaaa()
	// err := c.Visit(site + term + displayMode)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// results, _ := json.MarshalIndent(itens, "", " ")

	// _ = ioutil.WriteFile("results.json", results, 0644)

	data, err := ioutil.ReadFile("results.json")

	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(data, &itens)
	removeOutliers(&itens)
	for i := 0; i < len(itens)/2; i++ {
		for j := 1; j < len(itens)/2; j++ {
			similarity := strutil.Similarity(itens[i].Title, itens[j].Title, metrics.NewHamming())
			log.Println("Similarity ", itens[i].Title, itens[j].Title, similarity)
		}
	}

}

func removeOutliers(itens *[]Item) []Item {

	data := []float64{}
	for _, v := range *itens {
		data = append(data, v.Price)
	}

	std, _ := stats.StandardDeviation(data)
	mean, _ := stats.Mean(data)
	log.Println(std)
	q, _ := stats.Quartile(data)
	log.Println(q)
	cleanedData := []Item{}
	for _, v := range *itens {
		if math.Abs(v.Price) > math.Abs(mean-0.5*std) && math.Abs(v.Price) < math.Abs(mean+0.5*std) {
			cleanedData = append(cleanedData, v)
		}
	}

	log.Println(len(cleanedData))
	for _, v := range cleanedData {
		log.Println(v.Price)

	}
	return cleanedData

}
