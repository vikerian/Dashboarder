// services/web-server/main.go
package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"home-dashboard/shared/config"

	"github.com/gorilla/mux"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

type WebServer struct {
	config    *config.Config
	router    *mux.Router
	templates *template.Template
}

func NewWebServer(cfg *config.Config) (*WebServer, error) {
	tmpl, err := template.ParseFS(templates, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	ws := &WebServer{
		config:    cfg,
		router:    mux.NewRouter(),
		templates: tmpl,
	}

	ws.setupRoutes()
	return ws, nil
}

func (ws *WebServer) setupRoutes() {
	// Static files
	ws.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	// Pages
	ws.router.HandleFunc("/", ws.handleHome).Methods("GET")
	ws.router.HandleFunc("/traffic", ws.handleTraffic).Methods("GET")
	ws.router.HandleFunc("/sensors", ws.handleSensors).Methods("GET")
	ws.router.HandleFunc("/news", ws.handleNews).Methods("GET")
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Home Dashboard",
		"Time":       time.Now().Format("15:04:05"),
		"Date":       time.Now().Format("2006-01-02"),
		"APIBaseURL": fmt.Sprintf("http://localhost:%s/api", ws.config.APIPort),
	}

	if err := ws.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleTraffic(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Traffic Incidents - Ústecký kraj",
		"APIBaseURL": fmt.Sprintf("http://localhost:%s/api", ws.config.APIPort),
	}

	if err := ws.templates.ExecuteTemplate(w, "traffic.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleSensors(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "IoT Sensors",
		"APIBaseURL": fmt.Sprintf("http://localhost:%s/api", ws.config.APIPort),
	}

	if err := ws.templates.ExecuteTemplate(w, "sensors.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleNews(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "MMDecin News",
		"APIBaseURL": fmt.Sprintf("http://localhost:%s/api", ws.config.APIPort),
	}

	if err := ws.templates.ExecuteTemplate(w, "news.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) Start() error {
	addr := ":" + ws.config.WebPort
	log.Printf("Starting web server on %s", addr)

	return http.ListenAndServe(addr, ws.router)
}

func main() {
	cfg := config.Load()

	server, err := NewWebServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create web server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// services/web-server/templates/base.html
/*
<!DOCTYPE html>
<html lang="cs">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <nav class="navbar">
        <div class="container">
            <h1>Home Dashboard</h1>
            <ul class="nav-links">
                <li><a href="/">Overview</a></li>
                <li><a href="/traffic">Traffic</a></li>
                <li><a href="/sensors">Sensors</a></li>
                <li><a href="/news">News</a></li>
            </ul>
        </div>
    </nav>

    <main class="container">
        {{template "content" .}}
    </main>

    <footer>
        <p>Last updated: <span id="last-update">{{.Time}}</span></p>
    </footer>

    <script src="/static/app.js"></script>
</body>
</html>
*/

// services/web-server/templates/index.html
/*
{{template "base.html" .}}

{{define "content"}}
<div class="dashboard">
    <h2>System Overview</h2>

    <div class="stats-grid">
        <div class="stat-card">
            <h3>MMDecin Articles</h3>
            <div class="stat-value" id="article-count">-</div>
        </div>

        <div class="stat-card">
            <h3>Active Traffic Incidents</h3>
            <div class="stat-value" id="traffic-count">-</div>
        </div>

        <div class="stat-card">
            <h3>Connected Sensors</h3>
            <div class="stat-value" id="sensor-count">-</div>
        </div>

        <div class="stat-card">
            <h3>System Time</h3>
            <div class="stat-value" id="system-time">{{.Time}}</div>
        </div>
    </div>

    <div class="recent-section">
        <h3>Recent Activity</h3>
        <div id="recent-activity" class="activity-list">
            <p>Loading...</p>
        </div>
    </div>
</div>

<script>
    const API_BASE = '{{.APIBaseURL}}';

    async function loadStats() {
        try {
            const response = await fetch(`${API_BASE}/stats`);
            const data = await response.json();

            document.getElementById('article-count').textContent = data.mmdecin_articles || 0;
            document.getElementById('traffic-count').textContent = data.traffic?.active || 0;
            document.getElementById('sensor-count').textContent = data.sensors || 0;
        } catch (error) {
            console.error('Failed to load stats:', error);
        }
    }

    function updateTime() {
        const now = new Date();
        document.getElementById('system-time').textContent = now.toLocaleTimeString('cs-CZ');
        document.getElementById('last-update').textContent = now.toLocaleTimeString('cs-CZ');
    }

    // Load data on page load
    loadStats();

    // Update every 30 seconds
    setInterval(loadStats, 30000);
    setInterval(updateTime, 1000);
</script>
{{end}}
*/

// services/web-server/templates/traffic.html
/*
{{template "base.html" .}}

{{define "content"}}
<div class="traffic-page">
    <h2>Traffic Incidents - Ústecký kraj</h2>

    <div class="filter-bar">
        <button onclick="loadTraffic('active')" class="btn btn-primary">Active Only</button>
        <button onclick="loadTraffic('all')" class="btn">All Incidents</button>
    </div>

    <div id="traffic-list" class="incident-list">
        <p>Loading traffic data...</p>
    </div>
</div>

<script>
    const API_BASE = '{{.APIBaseURL}}';

    async function loadTraffic(filter = 'active') {
        const endpoint = filter === 'active' ? '/traffic/active' : '/traffic';

        try {
            const response = await fetch(`${API_BASE}${endpoint}`);
            const incidents = await response.json();

            const listEl = document.getElementById('traffic-list');

            if (incidents.length === 0) {
                listEl.innerHTML = '<p>No traffic incidents found.</p>';
                return;
            }

            listEl.innerHTML = incidents.map(incident => `
                <div class="incident-card ${incident.severity}">
                    <h3>${incident.type}</h3>
                    <p class="location">${incident.location}</p>
                    <p class="description">${incident.description}</p>
                    <p class="time">${new Date(incident.timestamp).toLocaleString('cs-CZ')}</p>
                </div>
            `).join('');
        } catch (error) {
            console.error('Failed to load traffic:', error);
            document.getElementById('traffic-list').innerHTML = '<p>Error loading traffic data.</p>';
        }
    }

    // Load active incidents on page load
    loadTraffic('active');
</script>
{{end}}
*/

// services/web-server/static/style.css
/*
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: #333;
    background-color: #f5f5f5;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

.navbar {
    background-color: #2c3e50;
    color: white;
    padding: 1rem 0;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.navbar h1 {
    display: inline-block;
    font-size: 1.5rem;
}

.nav-links {
    list-style: none;
    display: inline-block;
    float: right;
}

.nav-links li {
    display: inline-block;
    margin-left: 2rem;
}

.nav-links a {
    color: white;
    text-decoration: none;
    transition: opacity 0.3s;
}

.nav-links a:hover {
    opacity: 0.8;
}

main {
    padding: 2rem 0;
    min-height: calc(100vh - 200px);
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
}

.stat-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    text-align: center;
}

.stat-card h3 {
    font-size: 1rem;
    color: #666;
    margin-bottom: 0.5rem;
}

.stat-value {
    font-size: 2.5rem;
    font-weight: bold;
    color: #2c3e50;
}

.incident-card {
    background: white;
    padding: 1.5rem;
    margin-bottom: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    border-left: 4px solid #3498db;
}

.incident-card.high {
    border-left-color: #e74c3c;
}

.incident-card.medium {
    border-left-color: #f39c12;
}

.incident-card h3 {
    margin-bottom: 0.5rem;
    color: #2c3e50;
}

.incident-card .location {
    font-weight: bold;
    color: #555;
}

.incident-card .time {
    font-size: 0.9rem;
    color: #888;
    margin-top: 0.5rem;
}

.btn {
    display: inline-block;
    padding: 0.5rem 1rem;
    background-color: #ecf0f1;
    color: #2c3e50;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    text-decoration: none;
    transition: background-color 0.3s;
}

.btn:hover {
    background-color: #bdc3c7;
}

.btn-primary {
    background-color: #3498db;
    color: white;
}

.btn-primary:hover {
    background-color: #2980b9;
}

.filter-bar {
    margin-bottom: 1.5rem;
}

.filter-bar button {
    margin-right: 0.5rem;
}

footer {
    background-color: #34495e;
    color: white;
    text-align: center;
    padding: 1rem 0;
    position: relative;
    bottom: 0;
    width: 100%;
}
*/
