package main

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

type KEV struct {
	Title           string          `json:"title"`
	CatalogVersion  string          `json:"catalogVersion"`
	DateReleased    string          `json:"dateReleased"`
	Count           int64           `json:"count"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Vulnerability struct {
	CveID             string `json:"cveID"`
	VendorProject     string `json:"vendorProject"`
	Product           string `json:"product"`
	VulnerabilityName string `json:"vulnerabilityName"`
	DateAdded         string `json:"dateAdded"`
	ShortDescription  string `json:"shortDescription"`
	RequiredAction    string `json:"requiredAction"`
	DueDate           string `json:"dueDate"`
	Notes             string `json:"notes"`
}

func main() {

	url := "https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var kevData KEV
	err = json.Unmarshal(body, &kevData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Printf("KEV data: %+v\n", kevData)

}
