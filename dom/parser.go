package dom

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	resty "github.com/go-resty/resty/v2"
)

type Metals struct {
	Gold      float64
	Silver    float64
	Platinum  float64
	Palladium float64
}

type Parser struct {
	client *resty.Client
}

func Init() *Parser {
	client := resty.New()

	client.SetTimeout(10 * time.Second)

	transport := &http.Transport{

		DialContext: (&net.Dialer{

			Timeout: 10 * time.Second,
		}).DialContext,

		TLSHandshakeTimeout: 10 * time.Second,
	}

	client.SetTransport(transport)

	return &Parser{client: client}
}

func getDom(body []byte) *goquery.Document {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Print(err)
		return nil
	}

	return doc
}

func getMetalsFromDom(dom *goquery.Document) Metals {
	var (
		m   Metals
		err error
	)

	if m.Gold, err = strconv.ParseFloat(
		strings.ReplaceAll(
			dom.Find("div[data-room=RY0306000000XAU]").Nodes[0].FirstChild.FirstChild.Data,
			",", ""), 32); err != nil {
		return Metals{}
	}

	if m.Silver, err = strconv.ParseFloat(
		strings.ReplaceAll(
			dom.Find("div[data-room=RY0306000000XAG]").Nodes[0].FirstChild.FirstChild.Data,
			",", ""), 32); err != nil {
		return Metals{}
	}

	if m.Platinum, err = strconv.ParseFloat(
		strings.ReplaceAll(
			dom.Find("div[data-room=RY0306000000XPT]").Nodes[0].FirstChild.FirstChild.Data,
			",", ""), 32); err != nil {
		return Metals{}
	}

	if m.Palladium, err = strconv.ParseFloat(
		strings.ReplaceAll(
			dom.Find("div[data-room=RY0306000000XPD]").Nodes[0].FirstChild.FirstChild.Data,
			",", ""), 32); err != nil {
		return Metals{}
	}

	return m
}

func (p Parser) Pars(url string) Metals {
	res, err := p.client.R().Get(url)
	if err != nil {
		return Metals{}
	}

	if dom := getDom(res.Body()); dom != nil {
		return getMetalsFromDom(dom)
	}

	return Metals{}
}
