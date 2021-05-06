package airports

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// https://ourairports.com/help/data-dictionary.html
// "id","code","local_code","name","continent","iso_country","wikipedia_link","keywords"
// e.g. 306080,"US-CA","CA","California","NA","US","https://en.wikipedia.org/wiki/California",

const (
	colRegionID = iota
	colRegionCode
	colRegionLocalCode
	colRegionName
	colRegionContinent
	colRegionISOCountry
	colRegionWikipediaLink
	colRegionKeywords
)

type RegionData struct {
	ID            uint64
	ISOCode       string
	LocalCode     string
	Name          string
	Continent     string
	ISOCountry    string
	WikipediaLink string
	Keywords      string
}

type RegionDB struct {
	Regions map[string]*RegionData
}

func NewRegionDB() *RegionDB {
	return &RegionDB{
		Regions: make(map[string]*RegionData),
	}
}

func (db *RegionDB) Clear() {
	db.Regions = nil
}

func (db *RegionDB) Parse(file string, skipFirstLine bool) error {
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

		id, err := ParseUint(row[colRegionID])
		if err != nil {
			log.Println(line, err)
			continue
		}

		isoCode := row[colRegionCode]
		localCode := row[colRegionLocalCode]
		name := row[colRegionName]
		continent := row[colRegionContinent]
		isoContry := row[colRegionISOCountry]
		wikipedia := row[colRegionWikipediaLink]
		keywords := row[colRegionKeywords]

		region := &RegionData{
			ID:            id,
			ISOCode:       isoCode,
			LocalCode:     localCode,
			Name:          name,
			Continent:     continent,
			ISOCountry:    isoContry,
			WikipediaLink: wikipedia,
			Keywords:      keywords,
		}
		db.Regions[isoCode] = region
	}
}

func (db *RegionDB) FindByISOCode(isoCode string) *RegionData {
	if region, ok := db.Regions[isoCode]; ok {
		return region
	}
	return nil
}
