package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// this application can work as one of the agents' actuators
// based on the current intention, an agent can retrieve information about some topic
// sensors - beliefs (topic) - retrieve info from the web - sensor -> neural network - learn and improve an existent plan
// learn and improve an agent strategy, for instance

type Item struct {
	Title       string //colocar anotacao de len
	Description string
	Price       string
	User        string
	// UserRating  string
	Amount int
	New    bool
}

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

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
		item.Price = price
		item.User = user

		item.Amount, _ = strconv.Atoi(strings.Split(amount, " ")[0])

		item.New = strings.Contains(new, "Novo")
		log.Println(item)
	})

	// detailCollector.OnHTML("h1.ui-pdp-title", func(e *colly.HTMLElement) {
	// 	// el := e.Request.Visit(e.Attr("ol"))
	// 	// log.Println(el)
	// 	log.Println("Visiting product", e)

	// 	product := Item{}
	// 	product.Title = e.Text //ui-pdp-title

	// })

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	term := "gtx-1060"
	site := "https://lista.mercadolivre.com.br/"
	displayMode := "_DisplayType_LF"

	err := c.Visit(site + term + displayMode)
	if err != nil {
		fmt.Println(err)
	}

}
