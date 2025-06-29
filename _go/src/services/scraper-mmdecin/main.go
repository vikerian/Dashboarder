// services/scraper-mmdecin/main.go
package main

import (
	"log"
	"strings"
	"time"

	"home-dashboard/shared/config"
	"home-dashboard/shared/models"
	"home-dashboard/shared/mqtt"

	"github.com/gocolly/colly/v2"
)

const (
	scraperInterval = 30 * time.Minute
	mqttTopic       = "scrapers/mmdecin"
	baseURL         = "https://www.mmdecin.cz"
)

type Scraper struct {
	mqttClient *mqtt.Client
	config     *config.Config
	collector  *colly.Collector
}

func NewScraper(cfg *config.Config) (*Scraper, error) {
	mqttConfig := &mqtt.Config{
		Broker:   cfg.MQTTBroker,
		ClientID: "mmdecin-scraper",
	}

	mqttClient, err := mqtt.NewClient(mqttConfig)
	if err != nil {
		return nil, err
	}

	// Create a new collector with settings
	c := colly.NewCollector(
		colly.AllowedDomains("www.mmdecin.cz", "mmdecin.cz"),
		colly.MaxDepth(2),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// Set limits
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*mmdecin.cz*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	return &Scraper{
		mqttClient: mqttClient,
		config:     cfg,
		collector:  c,
	}, nil
}

func (s *Scraper) setupCollector() {
	// Reset callbacks for fresh start
	s.collector = s.collector.Clone()

	// Track visited articles to avoid duplicates
	visitedURLs := make(map[string]bool)

	// Handle main page and article listings
	s.collector.OnHTML("article, .news-item, .article, .post, .entry", func(e *colly.HTMLElement) {
		// Try multiple selectors for title
		title := strings.TrimSpace(e.ChildText("h1, h2, h3, .title, .headline, .post-title"))
		if title == "" {
			title = strings.TrimSpace(e.ChildText("a[href]"))
		}

		// Try multiple selectors for content/description
		content := strings.TrimSpace(e.ChildText("p, .content, .description, .excerpt, .summary"))
		if content == "" {
			// Try to get first paragraph
			content = strings.TrimSpace(e.ChildText("p:first-child"))
		}

		// Extract URL
		link := e.ChildAttr("a[href]", "href")
		if link == "" {
			// Try to find link in heading
			link = e.ChildAttr("h1 a, h2 a, h3 a", "href")
		}

		// Make relative URLs absolute
		absoluteURL := e.Request.AbsoluteURL(link)

		// Skip if we already processed this URL or if data is incomplete
		if title == "" || visitedURLs[absoluteURL] {
			return
		}

		visitedURLs[absoluteURL] = true

		// Truncate content if too long
		if len(content) > 500 {
			content = content[:497] + "..."
		}

		data := models.MMDecinData{
			Title:     title,
			Content:   content,
			URL:       absoluteURL,
			Timestamp: time.Now(),
		}

		s.publishArticle(data)

		// Visit the full article page for more details
		if link != "" && !strings.HasPrefix(link, "#") {
			e.Request.Visit(absoluteURL)
		}
	})

	// Look for pagination and follow next pages
	s.collector.OnHTML("a.next, a.pagination-next, .pagination a, .pager a", func(e *colly.HTMLElement) {
		nextURL := e.Attr("href")
		if nextURL != "" {
			e.Request.Visit(e.Request.AbsoluteURL(nextURL))
		}
	})

	// Handle news section links
	s.collector.OnHTML("a[href*='aktuality'], a[href*='zpravy'], a[href*='clanky'], a[href*='news']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" && !visitedURLs[e.Request.AbsoluteURL(link)] {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Handle individual article pages
	s.collector.OnHTML("body", func(e *colly.HTMLElement) {
		// Skip if this is a listing page
		if e.DOM.Find("article, .news-item").Length() > 1 {
			return
		}

		// Try to extract article from single article page
		title := strings.TrimSpace(e.ChildText("h1, .article-title, .entry-title"))

		// Try multiple content selectors
		var content string
		contentSelectors := []string{
			".article-content",
			".entry-content",
			".post-content",
			".content",
			"main article",
			"[itemprop='articleBody']",
		}

		for _, selector := range contentSelectors {
			content = strings.TrimSpace(e.ChildText(selector))
			if content != "" {
				break
			}
		}

		url := e.Request.URL.String()

		if title != "" && content != "" && !visitedURLs[url] {
			visitedURLs[url] = true

			// Clean content
			content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
			content = strings.ReplaceAll(content, "\t", " ")

			if len(content) > 1000 {
				content = content[:997] + "..."
			}

			data := models.MMDecinData{
				Title:     title,
				Content:   content,
				URL:       url,
				Timestamp: time.Now(),
			}

			s.publishArticle(data)
		}
	})

	// Set up error handling
	s.collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with error: %v", r.Request.URL, err)
	})

	// Log visited pages
	s.collector.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s", r.URL.String())
	})
}

func (s *Scraper) publishArticle(data models.MMDecinData) {
	message := models.MQTTMessage{
		Type:      "mmdecin_data",
		Timestamp: time.Now(),
		Data:      data,
	}

	if err := s.mqttClient.Publish(mqttTopic, message); err != nil {
		log.Printf("Failed to publish to MQTT: %v", err)
	} else {
		log.Printf("Published article: %s", data.Title)
	}
}

func (s *Scraper) scrapeMMDecin() error {
	// Set up collector with fresh callbacks
	s.setupCollector()

	// Start scraping from main page and common news sections
	startURLs := []string{
		baseURL,
		baseURL + "/aktuality",
		baseURL + "/zpravy",
		baseURL + "/novinky",
		baseURL + "/clanky",
	}

	for _, url := range startURLs {
		log.Printf("Starting scrape from: %s", url)
		if err := s.collector.Visit(url); err != nil {
			log.Printf("Failed to visit %s: %v", url, err)
		}
	}

	// Wait for collector to finish
	s.collector.Wait()

	return nil
}

func (s *Scraper) Start() {
	log.Println("Starting MMDecin scraper...")

	// Initial scrape
	if err := s.scrapeMMDecin(); err != nil {
		log.Printf("Initial scrape error: %v", err)
	}

	// Set up periodic scraping
	ticker := time.NewTicker(scraperInterval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running scheduled scrape...")
		if err := s.scrapeMMDecin(); err != nil {
			log.Printf("Scrape error: %v", err)
		}
	}
}

func main() {
	cfg := config.Load()

	scraper, err := NewScraper(cfg)
	if err != nil {
		log.Fatalf("Failed to create scraper: %v", err)
	}
	defer scraper.mqttClient.Disconnect()

	scraper.Start()
}
