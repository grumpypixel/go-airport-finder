package alphafoxtrot

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"sort"
)

// https://ourairports.com/help/data-dictionary.html
// header: "id","ident","type","name","latitude_deg","longitude_deg","elevation_ft","continent","iso_country","iso_region","municipality","scheduled_service","gps_code","iata_code","local_code","home_link","wikipedia_link","keywords"
// e.g. 3632,"KLAX","large_airport","Los Angeles International Airport",33.942501,-118.407997,125,"NA","US","US-CA","Los Angeles","yes","KLAX","LAX","LAX","https://www.flylax.com/","https://en.wikipedia.org/wiki/Los_Angeles_International_Airport",

const (
	colAirportID = iota
	colAirportIdent
	colAirportType
	colAirportName
	colAirportLatitudeDeg
	colAirportLongitudeDeg
	colAirportElevationFt
	colAirportContinent
	colAirportISOCountry
	colAirportISORegion
	colAirportMunicipality
	colAirportScheduledService
	colAirportGPSCode
	colAirportIATACode
	colAirportLocalCode
	colAirportHomeLink
	colAirportWikipediaLink
	colAirportKeywords
)

type AirportData struct {
	ID               uint64
	ICAOCode         string
	Type             string
	TypeFlag         uint64
	Name             string
	LatitudeDeg      float64
	LongitudeDeg     float64
	ElevationFt      int64
	Continent        string
	ISOCountry       string
	ISORegion        string
	Municipality     string
	ScheduledService bool
	GPSCode          string
	IATACode         string
	LocalCode        string
	HomeLink         string
	WikipediaLink    string
	Keywords         string
}

type AirportDB struct {
	Airports []*AirportData
}

func NewAirportDB() *AirportDB {
	return &AirportDB{
		Airports: make([]*AirportData, 0),
	}
}

func (db *AirportDB) Clear() {
	db.Airports = nil
}

func (db *AirportDB) Parse(file string, airportTypeFilter uint64, skipFirstLine bool) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	line := -1
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		line++

		if line == 0 && skipFirstLine {
			continue
		}

		typ := row[colAirportType]
		typeFlag := AirportTypeFromString(typ)
		if typeFlag == AirportTypeUnknown || typeFlag&airportTypeFilter == 0 {
			continue
		}

		id, err := ParseUint(row[colAirportID])
		if err != nil {
			log.Println(line, err)
			continue
		}

		icaoCode := row[colAirportIdent]
		name := row[colAirportName]

		latitude, err := ParseFloat(row[colAirportLatitudeDeg])
		if err != nil {
			log.Println(line, icaoCode, name, err)
			continue
		}

		longitude, err := ParseFloat(row[colAirportLongitudeDeg])
		if err != nil {
			log.Println(line, icaoCode, name, err)
			continue
		}

		elevation, _ := ParseInt(row[colAirportElevationFt])
		continent := row[colAirportContinent]
		country := row[colAirportISOCountry]
		region := row[colAirportISORegion]
		municipality := row[colAirportMunicipality]
		scheduledService := ParseBool(row[colAirportScheduledService])
		gpsCode := row[colAirportGPSCode]
		iataCode := row[colAirportIATACode]
		localCode := row[colAirportLocalCode]
		homeLink := row[colAirportHomeLink]
		wikipediaLink := row[colAirportWikipediaLink]
		keywords := row[colAirportKeywords]

		airport := &AirportData{
			ID:               id,
			ICAOCode:         icaoCode,
			Type:             typ,
			TypeFlag:         typeFlag,
			Name:             name,
			LatitudeDeg:      latitude,
			LongitudeDeg:     longitude,
			ElevationFt:      elevation,
			Continent:        continent,
			ISOCountry:       country,
			ISORegion:        region,
			Municipality:     municipality,
			ScheduledService: scheduledService,
			GPSCode:          gpsCode,
			IATACode:         iataCode,
			LocalCode:        localCode,
			HomeLink:         homeLink,
			WikipediaLink:    wikipediaLink,
			Keywords:         keywords,
		}
		db.Airports = append(db.Airports, airport)
	}
}

func (db *AirportDB) FindByAirportType(airportTypeFilter uint64) []*AirportData {
	airports := make([]*AirportData, 0)
	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		airports = append(airports, airport)
	}
	return airports
}

func (db *AirportDB) FindByICAOCode(icaoCode string) *AirportData {
	for _, airport := range db.Airports {
		if airport.ICAOCode == icaoCode {
			return airport
		}
	}
	return nil
}

func (db *AirportDB) FindByIATACode(iataCode string) *AirportData {
	for _, airport := range db.Airports {
		if airport.IATACode == iataCode {
			return airport
		}
	}
	return nil
}

func (db *AirportDB) FindByRegion(isoRegion string, airportTypeFilter uint64) []*AirportData {
	airports := make([]*AirportData, 0)
	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		if airport.ISORegion == isoRegion {
			airports = append(airports, airport)
		}
	}
	return airports
}

func (db *AirportDB) FindByCountry(isoCountry string, airportTypeFilter uint64) []*AirportData {
	airports := make([]*AirportData, 0)
	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		if airport.ISOCountry == isoCountry {
			airports = append(airports, airport)
		}
	}
	return airports
}

func (db *AirportDB) FindByContinent(continent string, airportTypeFilter uint64) []*AirportData {
	airports := make([]*AirportData, 0)
	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		if airport.Continent == continent {
			airports = append(airports, airport)
		}
	}
	return airports
}

func (db *AirportDB) FindNearestAirport(latitudeDeg, longitudeDeg, radius float64, airportTypeFilter uint64) *AirportData {
	var nearestAirport *AirportData = nil
	minDistance := math.MaxFloat64

	if radius < 0 {
		radius = math.MaxFloat64
	}

	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		distance := Distance(latitudeDeg, longitudeDeg, airport.LatitudeDeg, airport.LongitudeDeg)
		if distance < minDistance && distance <= radius {
			minDistance = distance
			nearestAirport = airport
		}
	}
	return nearestAirport
}

func (db *AirportDB) FindNearestAirports(latitudeDeg, longitudeDeg, radiusMeters float64, maxResults int, airportTypeFilter uint64) []*AirportData {
	if radiusMeters < 0 {
		radiusMeters = math.MaxFloat64
	}
	if maxResults < 0 {
		maxResults = math.MaxInt32
	}

	type Candidate struct {
		Airport  *AirportData
		Distance float64
	}

	candidates := make([]*Candidate, 0)
	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		distance := Distance(latitudeDeg, longitudeDeg, airport.LatitudeDeg, airport.LongitudeDeg)
		if distance <= radiusMeters {
			candidates = append(candidates, &Candidate{airport, distance})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Distance < candidates[j].Distance
	})

	airports := make([]*AirportData, 0, len(candidates))
	count := MinInt(len(candidates), maxResults)
	for i := 0; i < count; i++ {
		airports = append(airports, candidates[i].Airport)
	}
	return airports
}

func (db *AirportDB) FindAll(isoRegionFilter string, isoCountryFilter string, continentFilter string, airportTypeFilter uint64) []*AirportData {
	filterRegion := len(isoRegionFilter) > 0
	filterCountry := len(isoCountryFilter) > 0
	filterContinent := len(continentFilter) > 0

	airports := make([]*AirportData, 0)
	for _, airport := range db.Airports {
		if airport.TypeFlag&airportTypeFilter == 0 {
			continue
		}
		if filterRegion && airport.ISORegion != isoRegionFilter {
			continue
		}
		if filterCountry && airport.ISOCountry != isoCountryFilter {
			continue
		}
		if filterContinent && airport.Continent != continentFilter {
			continue
		}
		airports = append(airports, airport)
	}
	return airports
}

func FindNearestAirports(airports []*AirportData, latitudeDeg, longitudeDeg, radiusMeters float64, maxResults int) []*AirportData {
	if radiusMeters < 0 {
		radiusMeters = math.MaxFloat64
	}
	if maxResults < 0 {
		maxResults = math.MaxInt32
	}

	type Candidate struct {
		Airport  *AirportData
		Distance float64
	}

	candidates := make([]*Candidate, 0)
	for _, airport := range airports {
		distance := Distance(latitudeDeg, longitudeDeg, airport.LatitudeDeg, airport.LongitudeDeg)
		if distance <= radiusMeters {
			candidates = append(candidates, &Candidate{airport, distance})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Distance < candidates[j].Distance
	})

	result := make([]*AirportData, 0, len(candidates))
	count := MinInt(len(candidates), maxResults)
	for i := 0; i < count; i++ {
		result = append(result, candidates[i].Airport)
	}
	return result
}
