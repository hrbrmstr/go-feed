package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/feeds"
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
		log.Fatal("Error fetching data:", err) // https://pkg.go.dev/log#Fatal
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var kevData KEV
	err = json.Unmarshal(body, &kevData)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	resp.Body.Close()

	feed := &feeds.Feed{
		Title:       kevData.Title,
		Link:        &feeds.Link{Href: url},
		Description: "CISA's Known Exploited Vulnerability (KEV) catalog is an authoritative source of vulnerabilities that have been exploited in the wild, helping organizations prioritize remediation efforts to reduce cyber risks",
		Author:      &feeds.Author{Name: "CISA", Email: "Central@cisa.dhs.gov"},
		Created:     time.Now(),
	}

	slices.Reverse(kevData.Vulnerabilities) // https://pkg.go.dev/slices@master#Reverse

	for _, vulnerability := range kevData.Vulnerabilities {

		kevDate, _ := time.Parse("2006-01-02", vulnerability.DateAdded) // https://pkg.go.dev/time#Parse

		item := &feeds.Item{
			Title:       vulnerability.VulnerabilityName,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", strings.ToLower(vulnerability.CveID))},
			Description: vulnerability.ShortDescription,
			Id:          fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", strings.ToLower(vulnerability.CveID)),
			Created:     kevDate,
		}

		feed.Items = append(feed.Items, item)
		
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("kev-rss.xml")
	if err != nil {
		log.Fatal("Error:", err)
	}

	_, err = file.WriteString(rss)
	if err != nil {
		log.Fatal("Error:", err)
	}
	
	file.Close()

}
