package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Component struct {
	ID                       int     `json:"id"`
	Name                     string  `json:"name"`
	Type                     string  `json:"type"`
	Brand                    string  `json:"brand"`
	Model                    string  `json:"model"`
	Socket                   *string `json:"socket"`
	FormFactor               *string `json:"form_factor"`
	Wattage                  *int    `json:"wattage"`
	VRAM                     *int    `json:"vram"`
	RAMType                  *string `json:"ram_type"`
	RAMCapacity              *int    `json:"ram_capacity"`
	RAMSpeed                 *int    `json:"ram_speed"`
	Tier                     string  `json:"tier"`
	WorkloadScoreGaming      int     `json:"workload_score_gaming"`
	WorkloadScoreEditing     int     `json:"workload_score_editing"`
	WorkloadScoreProgramming int     `json:"workload_score_programming"`
	WorkloadScoreGeneral     int     `json:"workload_score_general"`
	Price                    float64 `json:"price"` // Live scraped price
}

type BuildRequest struct {
	Workload         string `json:"workload"`          // "gaming", "editing", "programming", "general"
	InCountry        bool   `json:"in_country"`
	PreferencesCooling string `json:"preferences_cooling"` // "air", "liquid"
	PreferencesGPU     string `json:"preferences_gpu"`     // "igpu", "discrete"
	Resolution       string `json:"resolution"`
	ColorDepth       string `json:"color_depth"`
	Budget           float64 `json:"budget"`
	LockedComponents []int  `json:"locked_components"`
}

type ScraperResponse struct {
	Query     string  `json:"query"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	InCountry bool    `json:"in_country"`
	Source    string  `json:"source"`
}

var db *sql.DB

func initDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "root"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1:3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "pcpartpicker"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = d.Ping(); err != nil {
		return nil, err
	}
	return d, nil
}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/api/components", getComponents)
	e.POST("/api/build", generateBuild)

	e.Logger.Fatal(e.Start(":8080"))
}

func getComponents(c echo.Context) error {
	rows, err := db.Query("SELECT id, name, type, brand, model, socket, form_factor, wattage, vram, ram_type, ram_capacity, ram_speed, tier, workload_score_gaming, workload_score_editing, workload_score_programming, workload_score_general FROM Components")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var comp Component
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.Type, &comp.Brand, &comp.Model, &comp.Socket, &comp.FormFactor, &comp.Wattage, &comp.VRAM, &comp.RAMType, &comp.RAMCapacity, &comp.RAMSpeed, &comp.Tier, &comp.WorkloadScoreGaming, &comp.WorkloadScoreEditing, &comp.WorkloadScoreProgramming, &comp.WorkloadScoreGeneral); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		components = append(components, comp)
	}

	return c.JSON(http.StatusOK, components)
}

var fetchPriceFunc = fetchPriceFromScraper // Allows mocking in tests

func fetchPriceFromScraper(componentName string, inCountry bool) (float64, error) {
	apiURL := fmt.Sprintf("http://localhost:8000/scrape?query=%s&in_country=%t", url.QueryEscape(componentName), inCountry)
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var sResp ScraperResponse
	if err := json.Unmarshal(body, &sResp); err != nil {
		return 0, err
	}

	return sResp.Price, nil
}

func generateBuild(c echo.Context) error {
	var req BuildRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Fetch all components
	rows, err := db.Query("SELECT id, name, type, socket, wattage, workload_score_gaming, workload_score_editing, workload_score_programming, workload_score_general FROM Components")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var allComponents []Component
	for rows.Next() {
		var comp Component
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.Type, &comp.Socket, &comp.Wattage, &comp.WorkloadScoreGaming, &comp.WorkloadScoreEditing, &comp.WorkloadScoreProgramming, &comp.WorkloadScoreGeneral); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		allComponents = append(allComponents, comp)
	}

	// Parallel fetch prices
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range allComponents {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			price, err := fetchPriceFunc(allComponents[i].Name, req.InCountry)
			if err == nil {
				mu.Lock()
				allComponents[i].Price = price
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	// Simple Build Engine logic:
	var finalBuild []Component
	var locked []Component

	// Filter and append locked components first
	for _, comp := range allComponents {
		for _, lockedID := range req.LockedComponents {
			if comp.ID == lockedID {
				locked = append(locked, comp)
				finalBuild = append(finalBuild, comp)
			}
		}
	}

	// Requirement Types to satisfy
	requiredTypes := []string{"CPU", "Motherboard", "GPU", "PSU"} // Simplified
	if req.PreferencesGPU == "igpu" {
		requiredTypes = []string{"CPU", "Motherboard", "PSU"}
	}

	currentSocket := ""
	currentWattage := 0

	for _, l := range locked {
		if l.Type == "CPU" || l.Type == "Motherboard" {
			if l.Socket != nil {
				currentSocket = *l.Socket
			}
		}
		if l.Wattage != nil && (l.Type == "CPU" || l.Type == "GPU") {
			currentWattage += *l.Wattage
		}
	}

	totalCost := 0.0
	for _, l := range locked {
		totalCost += l.Price
	}

	remainingBudget := req.Budget - totalCost

	// Adjust target specs
	// We'll give a score bump for high tier components if the user wants 4k or 10-bit
	highEndSpecs := false
	if req.Resolution == "4k" || req.ColorDepth == "10bit" {
		highEndSpecs = true
	}

	// Dynamic budgeting: iteratively downgrade if we can't afford all parts.
	// For simplicity, we implement a greedy selection but retry with relaxed tiers if budget fails.
	success := false
	for attempts := 0; attempts < 3 && !success; attempts++ {
		tempBuild := append([]Component(nil), finalBuild...) // copy
		tempRemaining := remainingBudget
		tempWattage := currentWattage
		tempSocket := currentSocket

		allTypesSatisfied := true

		for _, reqType := range requiredTypes {
			alreadyHas := false
			for _, b := range tempBuild {
				if b.Type == reqType {
					alreadyHas = true
					break
				}
			}
			if alreadyHas {
				continue
			}

			var bestComponent *Component
			var bestScore float64 = -1

			for i, comp := range allComponents {
				if comp.Type != reqType {
					continue
				}

				// Price check
				if comp.Price > tempRemaining {
					continue
				}

				// Attempt relaxation based on attempt count
				if attempts == 0 && comp.Tier != "S" && comp.Tier != "A" && highEndSpecs {
					// Initially try to get only good tiers for high end specs
					continue
				} else if attempts == 1 && comp.Tier == "D" {
					// Try to avoid worst tier
					continue
				}

				// Compatibility Checks
				if (reqType == "Motherboard" || reqType == "CPU") && tempSocket != "" {
					if comp.Socket != nil && *comp.Socket != tempSocket {
						continue
					}
				}

				if reqType == "PSU" && comp.Wattage != nil {
					// Needs 20% headroom
					if *comp.Wattage < int(float64(tempWattage)*1.2) {
						continue
					}
				}

				// Scoring based on workload
				score := 0
				switch strings.ToLower(req.Workload) {
				case "gaming":
					score = comp.WorkloadScoreGaming
				case "editing":
					score = comp.WorkloadScoreEditing
				case "programming":
					score = comp.WorkloadScoreProgramming
				default:
					score = comp.WorkloadScoreGeneral
				}

				// Calculate a value score (score per dollar) to maximize budget efficiency
				valueScore := float64(score)
				if comp.Price > 0 {
					valueScore = float64(score) / comp.Price
				}

				// We prioritize higher tier and better value within budget
				if valueScore > bestScore {
					bestScore = valueScore
					bestComponent = &allComponents[i]
				}
			}

			if bestComponent != nil {
				tempBuild = append(tempBuild, *bestComponent)
				tempRemaining -= bestComponent.Price

				if bestComponent.Type == "CPU" || bestComponent.Type == "Motherboard" {
					if bestComponent.Socket != nil {
						tempSocket = *bestComponent.Socket
					}
				}
				if bestComponent.Wattage != nil && (bestComponent.Type == "CPU" || bestComponent.Type == "GPU") {
					tempWattage += *bestComponent.Wattage
				}
			} else {
				allTypesSatisfied = false
				break // Failed to find a component for this type within constraints
			}
		}

		if allTypesSatisfied {
			finalBuild = tempBuild
			remainingBudget = tempRemaining
			success = true
		}
	}

	if !success {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Budget too low for requested parts or no compatible components found."})
	}

	totalCost = req.Budget - remainingBudget

	return c.JSON(http.StatusOK, map[string]interface{}{
		"build":      finalBuild,
		"total_cost": totalCost,
	})
}
