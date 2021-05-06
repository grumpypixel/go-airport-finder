package airports

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// https://ourairports.com/help/data-dictionary.html
// "id","code","name","continent","wikipedia_link","keywords"
// e.g. 302755,"US","United States","NA","https://en.wikipedia.org/wiki/United_States","America"

const (
	colCountryID = iota
	colCountryCode
	colCountryName
	colCountryContinent
	colCountryWikipediaLink
	colCountryKeywords
)

type CountryData struct {
	ID            uint64
	ISOCode       string
	Name          string
	Continent     string
	WikipediaLink string
	Keywords      string
}

type CountryDB struct {
	Countries map[string]*CountryData
}

func NewCountryDB() *CountryDB {
	return &CountryDB{
		Countries: make(map[string]*CountryData),
	}
}

func (db *CountryDB) Clear() {
	db.Countries = nil
}

func (db *CountryDB) Parse(file string, skipFirstLine bool) error {
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

		id, err := ParseUint(row[colCountryID])
		if err != nil {
			log.Println(line, err)
			continue
		}

		isoCode := row[colCountryCode]
		name := row[colCountryName]
		continent := row[colCountryContinent]
		wikipedia := row[colCountryWikipediaLink]
		keywords := row[colCountryKeywords]

		country := &CountryData{
			ID:            id,
			ISOCode:       isoCode,
			Name:          name,
			Continent:     continent,
			WikipediaLink: wikipedia,
			Keywords:      keywords,
		}
		db.Countries[isoCode] = country
	}
}

func (db *CountryDB) FindByISOCode(isoCode string) *CountryData {
	if country, ok := db.Countries[isoCode]; ok {
		return country
	}
	return nil
}
