package main

import (
	//"bytes"
	//"fmt"
	//"string"

	//"github.com/gocolly/colly"
	"time"

	"github.com/gocolly/colly"
)

type Scrapper struct {
	Collector *colly.Collector
	Scrapes   []struct {
		Url        string
		Scrapped   time.Time
		RawContent []byte
		FileName   string
	}
}

func NewScrapper() *Scrapper {
	return new(Scrapper)
}
