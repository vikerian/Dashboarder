// services/scraper-traffic/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"home-dashboard/shared/config"
	"home-dashboard/shared/models"
	"home-dashboard/shared/mqtt"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const (
	scraperInterval = 15 * time.Minute
	mqttTopic       = "scrapers/traffic"
)

type TrafficScraper struct {
	mqttClient *mqtt.Client
	config     *config.Config
	mu         sync.Mutex
	incidents  map[string]bool // Track published incidents to avoid duplicates
}

func NewTrafficScraper(cfg *config.Config) (*TrafficScraper, error) {
	mqttConfig := &mqtt.Config{
		Broker:   cfg.MQTTBroker,
		ClientID: "traffic-scraper",
	}

	mqttClient, err := mqtt.NewClient(mqttConfig)
	if err != nil {
		return nil, err
	}

	return &TrafficScraper{
		mqttClient: mqttClient,
		config:     cfg,
		incidents:  make(map[string]bool),
	}, nil
}

func (s *TrafficScraper) createCollector(allowedDomains ...string) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		colly.MaxDepth(2),
		colly.Async(true),
	)

	// Set realistic browser headers
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping %s: %v", r.Request.URL, err)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "cs-CZ,cs;q=0.9,en;q=0.8")
		log.Printf("Visiting: %s", r.URL.String())
	})

	return c
}

func (s *TrafficScraper) scrapeRSD() error {
	log.Println("Starting RSD scraper...")
	c := s.createCollector("www.rsd.cz", "rsd.cz")

	// Handle traffic events on RSD website
	c.OnHTML(".dopravni-info-udalost, .traffic-event, .incident", func(e *colly.HTMLElement) {
		s.parseRSDIncident(e)
	})

	// Also look for table rows with incidents
	c.OnHTML("tr.incident-row, tbody tr", func(e *colly.HTMLElement) {
		// Check if this row contains traffic data
		if e.ChildText("td") != "" {
			s.parseRSDTableRow(e)
		}
	})

	// Visit RSD URLs
	urls := []string{
		"https://www.rsd.cz/dopravni-info",
		"https://www.rsd.cz/wps/portal/web/dopravni-info/omezeni-provozu",
		"https://www.rsd.cz/wps/portal/web/dopravni-info/uzavirky",
	}

	for _, url := range urls {
		if err := c.Visit(url); err != nil {
			log.Printf("Failed to visit RSD %s: %v", url, err)
		}
	}

	c.Wait()
	return nil
}

func (s *TrafficScraper) parseRSDIncident(e *colly.HTMLElement) {
	location := strings.TrimSpace(e.ChildText(".misto, .location"))
	if location == "" {
		// Try alternative selectors
		location = strings.TrimSpace(e.ChildText("h3, h4"))
	}

	if !s.isUsteckyKraj(location) {
		return
	}

	description := strings.TrimSpace(e.ChildText(".popis, .description, p"))
	incidentType := s.detectIncidentType(e, description)
	severity := s.determineSeverity(e, description)

	// Extract dates if available
	dateFrom := e.ChildText(".od, .date-from")
	dateTo := e.ChildText(".do, .date-to")
	if dateFrom != "" || dateTo != "" {
		description = fmt.Sprintf("%s (Od: %s Do: %s)", description, dateFrom, dateTo)
	}

	incident := models.TrafficIncident{
		Type:        incidentType,
		Location:    location,
		Description: description,
		Severity:    severity,
		Timestamp:   time.Now(),
	}

	s.publishIncident(incident)
}

func (s *TrafficScraper) parseRSDTableRow(e *colly.HTMLElement) {
	// Extract data from table cells
	cells := e.ChildTexts("td")
	if len(cells) < 3 {
		return
	}

	location := strings.TrimSpace(cells[0])
	if !s.isUsteckyKraj(location) {
		return
	}

	incident := models.TrafficIncident{
		Type:        "Dopravní omezení",
		Location:    location,
		Description: strings.TrimSpace(cells[1]),
		Severity:    "normal",
		Timestamp:   time.Now(),
	}

	if len(cells) > 2 {
		// Additional info might be in 3rd column
		incident.Description += " " + strings.TrimSpace(cells[2])
	}

	s.publishIncident(incident)
}

func (s *TrafficScraper) scrapeDopravniInfo() error {
	log.Println("Starting DopravniInfo scraper...")
	c := s.createCollector("www.dopravniinfo.cz", "dopravniinfo.cz")

	// Main incident selector
	c.OnHTML(".event-item, .traffic-event, .incident", func(e *colly.HTMLElement) {
		location := strings.TrimSpace(e.ChildText(".location, .misto"))

		if !s.isUsteckyKraj(location) {
			return
		}

		incident := models.TrafficIncident{
			Type:        e.Attr("data-type"),
			Location:    location,
			Description: strings.TrimSpace(e.ChildText(".description, .popis")),
			Severity:    s.determineSeverity(e, e.Text),
			Timestamp:   time.Now(),
		}

		// Try to get coordinates
		if lat := e.Attr("data-lat"); lat != "" {
			// Parse coordinates - simplified for example
			incident.Lat = 50.6833
		}
		if lng := e.Attr("data-lng"); lng != "" {
			incident.Lng = 14.0333
		}

		s.publishIncident(incident)
	})

	// Handle map data in scripts
	c.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "mapData") || strings.Contains(e.Text, "incidents") {
			s.extractJSONIncidents(e.Text)
		}
	})

	urls := []string{
		"https://www.dopravniinfo.cz/ustecky-kraj",
		"https://www.dopravniinfo.cz/nehody",
		"https://www.dopravniinfo.cz/omezeni",
	}

	for _, url := range urls {
		if err := c.Visit(url); err != nil {
			log.Printf("Failed to visit DopravniInfo %s: %v", url, err)
		}
	}

	c.Wait()
	return nil
}

func (s *TrafficScraper) scrapeIDOSTraffic() error {
	log.Println("Starting IDOS traffic scraper...")

	// IDOS has traffic info API
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://idos.idnes.cz/vlakyautobusymhdvse/spojeni/vysledky/", nil)
	if err != nil {
		return err
	}

	// This would need proper API endpoint
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch IDOS: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// Parse response...
	return nil
}

func (s *TrafficScraper) extractJSONIncidents(scriptContent string) {
	// Look for JSON data in script
	jsonPattern := regexp.MustCompile(`(?s)(?:incidents|events|data)\s*=\s*(\[.*?\]);`)
	matches := jsonPattern.FindStringSubmatch(scriptContent)

	if len(matches) > 1 {
		var incidents []map[string]interface{}
		if err := json.Unmarshal([]byte(matches[1]), &incidents); err == nil {
			for _, inc := range incidents {
				if location, ok := inc["location"].(string); ok && s.isUsteckyKraj(location) {
					incident := models.TrafficIncident{
						Type:        fmt.Sprintf("%v", inc["type"]),
						Location:    location,
						Description: fmt.Sprintf("%v", inc["description"]),
						Severity:    "normal",
						Timestamp:   time.Now(),
					}

					if lat, ok := inc["lat"].(float64); ok {
						incident.Lat = lat
					}
					if lng, ok := inc["lng"].(float64); ok {
						incident.Lng = lng
					}

					s.publishIncident(incident)
				}
			}
		}
	}
}

func (s *TrafficScraper) detectIncidentType(e *colly.HTMLElement, text string) string {
	classes := e.Attr("class")
	textLower := strings.ToLower(text)

	if strings.Contains(classes, "nehoda") || strings.Contains(textLower, "nehoda") {
		return "Nehoda"
	} else if strings.Contains(classes, "uzavirka") || strings.Contains(textLower, "uzavírka") {
		return "Uzavírka"
	} else if strings.Contains(classes, "omezeni") || strings.Contains(textLower, "omezení") {
		return "Omezení provozu"
	} else if strings.Contains(textLower, "oprava") || strings.Contains(textLower, "údržba") {
		return "Údržba silnice"
	} else if strings.Contains(textLower, "kolona") {
		return "Kolona"
	}

	return "Dopravní událost"
}

func (s *TrafficScraper) determineSeverity(e *colly.HTMLElement, text string) string {
	classes := e.Attr("class")
	textLower := strings.ToLower(text)

	if strings.Contains(classes, "critical") || strings.Contains(classes, "high") ||
		strings.Contains(textLower, "uzavírka") || strings.Contains(textLower, "neprůjezdná") {
		return "high"
	} else if strings.Contains(classes, "medium") ||
		strings.Contains(textLower, "omezení") || strings.Contains(textLower, "zúžení") {
		return "medium"
	}

	return "normal"
}

func (s *TrafficScraper) isUsteckyKraj(location string) bool {
	if location == "" {
		return false
	}

	locationLower := strings.ToLower(location)

	// Cities and towns in Ústecký kraj
	keywords := []string{
		"ústí", "usti", "ústec", "ustec",
		"děčín", "decin",
		"teplice",
		"most",
		"chomutov",
		"litoměřice", "litomerice",
		"louny",
		"žatec", "zatec",
		"roudnice",
		"lovosice",
		"bílina", "bilina",
		"duchcov",
		"klášterec", "klasterec",
		"kadaň", "kadan",
		"litvínov", "litvinov",
		"jirkov",
		"varnsdorf",
		"rumburk",
	}

	for _, keyword := range keywords {
		if strings.Contains(locationLower, keyword) {
			return true
		}
	}

	// Check road numbers typical for the region
	roadPattern := regexp.MustCompile(`\b(I/[0-9]+|II/[0-9]+|D8|E442|E55)\b`)
	if matches := roadPattern.FindAllString(location, -1); len(matches) > 0 {
		usteckyRoads := []string{
			"D8",   // Main highway through the region
			"I/13", // Děčín - Liberec
			"I/15", // Most - Litvínov
			"I/27", // Most - Plzeň
			"I/30", // Ústí nad Labem - Lovosice
			"I/62", // Děčín - Česká Lípa
			"I/63", // Ústí nad Labem
			"II/240", "II/247", "II/253", "II/254", "II/255",
			"II/257", "II/258", "II/260", "II/261", "II/262",
			"II/263", "II/264", "II/265", "II/266",
		}

		for _, match := range matches {
			for _, road := range usteckyRoads {
				if match == road {
					return true
				}
			}
		}
	}

	return false
}

func (s *TrafficScraper) publishIncident(incident models.TrafficIncident) {
	// Create unique key to avoid duplicates
	key := fmt.Sprintf("%s|%s|%s", incident.Location, incident.Type, incident.Description)

	s.mu.Lock()
	if s.incidents[key] {
		s.mu.Unlock()
		return
	}
	s.incidents[key] = true
	s.mu.Unlock()

	message := models.MQTTMessage{
		Type:      "traffic_incident",
		Timestamp: time.Now(),
		Data:      incident,
	}

	if err := s.mqttClient.Publish(mqttTopic, message); err != nil {
		log.Printf("Failed to publish incident: %v", err)
	} else {
		log.Printf("Published traffic incident: %s at %s", incident.Type, incident.Location)
	}
}

func (s *TrafficScraper) clearOldIncidents() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear the incidents map periodically to allow for updates
	s.incidents = make(map[string]bool)
	log.Println("Cleared incident cache")
}

func (s *TrafficScraper) Start() {
	log.Println("Starting traffic scraper...")

	// Initial scrape
	s.runScraping()

	// Set up periodic scraping
	scraperTicker := time.NewTicker(scraperInterval)
	defer scraperTicker.Stop()

	// Clear cache every hour
	clearTicker := time.NewTicker(1 * time.Hour)
	defer clearTicker.Stop()

	for {
		select {
		case <-scraperTicker.C:
			log.Println("Running scheduled traffic scrape...")
			s.runScraping()
		case <-clearTicker.C:
			s.clearOldIncidents()
		}
	}
}

func (s *TrafficScraper) runScraping() {
	var wg sync.WaitGroup

	scrapers := []func() error{
		s.scrapeRSD,
		s.scrapeDopravniInfo,
		s.scrapeIDOSTraffic,
	}

	for _, scraper := range scrapers {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				log.Printf("Scraper error: %v", err)
			}
		}(scraper)
	}

	wg.Wait()
	log.Printf("Scraping complete. Total incidents tracked: %d", len(s.incidents))
}

func main() {
	cfg := config.Load()

	scraper, err := NewTrafficScraper(cfg)
	if err != nil {
		log.Fatalf("Failed to create traffic scraper: %v", err)
	}
	defer scraper.mqttClient.Disconnect()

	scraper.Start()
}
