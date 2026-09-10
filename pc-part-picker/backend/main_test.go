package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/labstack/echo/v4"
)

import "os"

// A simple test ensuring the endpoint handles locking and returns correct data.
func TestGenerateBuildAPI(t *testing.T) {
	// Mock the external scraper API call
	fetchPriceFunc = func(componentName string, inCountry bool) (float64, error) {
		if !inCountry {
			return 120.0, nil
		}
		return 100.0, nil
	}

	// Need to initialize db connection for testing
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "root")
	os.Setenv("DB_HOST", "127.0.0.1:3306")
	os.Setenv("DB_NAME", "pcpartpicker")

	var err error
	db, err = initDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	reqBody := BuildRequest{
		Workload:         "gaming",
		InCountry:        true,
		PreferencesCooling: "air",
		PreferencesGPU:     "discrete",
		Budget:           2000.0, // generous budget
		LockedComponents: []int{1}, // Lock AMD Ryzen 5 7600
	}

	jsonBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/build", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := generateBuild(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	build, ok := resp["build"].([]interface{})
	if !ok {
		t.Fatalf("Expected build array in response")
	}

	if len(build) == 0 {
		t.Errorf("Expected build to not be empty")
	}

	// Verify locked component is in build
	foundLocked := false
	for _, compInterface := range build {
		compMap := compInterface.(map[string]interface{})
		if int(compMap["id"].(float64)) == 1 {
			foundLocked = true
			break
		}
	}

	if !foundLocked {
		t.Errorf("Expected locked component (ID: 1) to be in the final build")
	}
}

func TestGenerateBuildAPI_CrossBorder(t *testing.T) {
	// Mock the external scraper API call
	fetchPriceFunc = func(componentName string, inCountry bool) (float64, error) {
		if !inCountry {
			return 120.0, nil
		}
		return 100.0, nil
	}

	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "root")
	os.Setenv("DB_HOST", "127.0.0.1:3306")
	os.Setenv("DB_NAME", "pcpartpicker")
	var err error
	db, err = initDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	reqBody := BuildRequest{
		Workload:         "general",
		InCountry:        false,
		PreferencesCooling: "air",
		PreferencesGPU:     "igpu",
		Budget:           2000.0, // increased budget so it doesn't fail the new stricter budget logic
		LockedComponents: []int{},
	}

	jsonBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/build", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := generateBuild(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	totalCost, ok := resp["total_cost"].(float64)
	if !ok {
		t.Fatalf("Expected total_cost in response")
	}

	if totalCost <= 0 {
		t.Errorf("Expected total cost to be greater than 0")
	}
}
