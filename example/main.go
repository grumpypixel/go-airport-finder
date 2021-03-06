package main

import (
	"fmt"
	"log"
	"os"
	"path"

	alphafoxtrot "github.com/grumpypixel/go-airport-finder"
)

func main() {
	// Create data directory if it doesn't exist
	// and download airport data from OurAirports.com
	dataDir := "./data"
	validateData(dataDir)

	// Create an airport finder instance
	finder := alphafoxtrot.NewAirportFinder()

	// LoadOptions come with preset filepaths
	options := alphafoxtrot.PresetLoadOptions(dataDir)

	// specify a filter what airport types should be loaded (see def.go)
	// to only load heliports and seaplane bases, you set the filter to be: filter := AirportTypeHeliport | AirportTypeSeaplaneBase
	filter := alphafoxtrot.AirportTypeAll

	// Load the data into memory
	if err := finder.Load(options, filter); len(err) > 0 {
		log.Println("errors:", err)
	}

	// Find an airport by its ICAO code
	airport := finder.FindAirportByICAOCode("KLAX")
	if airport != nil {
		fmt.Println("Airport by ICAO code:", *airport)
	}

	// Find an airport by its IATA code
	airport = finder.FindAirportByIATACode("DUS")
	if airport != nil {
		fmt.Println("\nAirport by IATA code:", *airport)
	}

	// Find the nearest active airport within a given radius
	radiusInMeters := alphafoxtrot.KilometersToMeters(25)
	airport = finder.FindNearestAirport(33.942501, -118.407997, radiusInMeters, alphafoxtrot.AirportTypeActive)
	if airport != nil {
		fmt.Println("\nNearest airport:", *airport)
	}

	// Find active airports within a given radius
	radiusInMeters = alphafoxtrot.NauticalMilesToMeters(100.0)
	maxResults := 10
	airports := finder.FindNearestAirports(33.942501, -118.407997, radiusInMeters, maxResults, alphafoxtrot.AirportTypeRunways)
	if len(airports) > 0 {
		fmt.Println("\nNearest airports:")
		for i, airport := range airports {
			fmt.Printf("#%d: %v\n", i+1, *airport)
		}
	}

	// Find all large airports in a region
	radiusInMeters = alphafoxtrot.MilesToMeters(50.0)
	maxResults = 10
	airports = finder.FindNearestAirportsByRegion("US-CA", 33.942501, -118.407997, radiusInMeters, maxResults, alphafoxtrot.AirportTypeLarge)
	if len(airports) > 0 {
		fmt.Println("\nNearest Airports by region:")
		for i, airport := range airports {
			fmt.Printf("#%d: %v\n", i+1, *airport)
		}
	}

	// Find all large and medium airports in a country
	radiusInMeters = alphafoxtrot.NauticalMilesToMeters(50.0)
	maxResults = 10
	airports = finder.FindNearestAirportsByCountry("US", 33.942501, -118.407997, radiusInMeters, maxResults, alphafoxtrot.AirportTypeMedium)
	if len(airports) > 0 {
		fmt.Println("\nNearest Airports by Country:")
		for i, airport := range airports {
			fmt.Printf("#%d: %v\n", i+1, *airport)
		}
	}

	// Find the nearest navaids within a given radius
	radiusInMeters = alphafoxtrot.KilometersToMeters(50.0)
	navaids := finder.FindNearestNavaids(33.942501, -118.407997, radiusInMeters, maxResults)
	if len(navaids) > 0 {
		fmt.Println("\nNearest navaids:")
		for i, navaid := range navaids {
			fmt.Printf("#%d %v\n", i+1, *navaid)
		}
	}

	// Find the navaids associated with an airport (by ICAO code)
	airportICAOCode := "CYYC"
	navaids = finder.FindNavaidsByAirportICAOCode(airportICAOCode)
	if len(navaids) > 0 {
		fmt.Println("\nNavaids associated with airport ICAO code:")
		for i, navaid := range navaids {
			fmt.Printf("#%d %v\n", i+1, *navaid)
		}
	}
}

func validateData(dataDir string) {
	downloadFiles := false
	for _, filename := range alphafoxtrot.OurAirportsFiles {
		filepath := path.Join(dataDir, filename)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			downloadFiles = true
			break
		}
	}
	if downloadFiles {
		fmt.Println("Downloading CSV files from OurAirports.com...")
		alphafoxtrot.DownloadDatabase(dataDir)
	}
}
