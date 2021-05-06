package airports

import (
	"fmt"
	"path/filepath"
)

type AirportFinder struct {
	airportDB   *AirportDB
	frequencyDB *FrequencyDB
	runwayDB    *RunwayDB
	regionDB    *RegionDB
	countryDB   *CountryDB
	navaidDB    *NavaidDB
}

type LoadOptions struct {
	AirportsFilename    string // required, usually: airports.csv
	FrequenciesFilename string // optional, usually: airport-frequencies.csv
	RunwaysFilename     string // optional, usually: runways.csv
	RegionsFilename     string // optional, usually: regions.csv
	CountriesFilename   string // optional, usually: countries.csv
	NavaidsFilename     string // optional, usually: navaids.csv
}

func PresetLoadOptions(baseDir string) *LoadOptions {
	return &LoadOptions{
		AirportsFilename:    filepath.Join(baseDir, OurAirportsFiles[AirportsFileKey]),
		FrequenciesFilename: filepath.Join(baseDir, OurAirportsFiles[FrequenciesFileKey]),
		RunwaysFilename:     filepath.Join(baseDir, OurAirportsFiles[RunwaysFileKey]),
		RegionsFilename:     filepath.Join(baseDir, OurAirportsFiles[RegionsFileKey]),
		CountriesFilename:   filepath.Join(baseDir, OurAirportsFiles[CountriesFileKey]),
		NavaidsFilename:     filepath.Join(baseDir, OurAirportsFiles[NavaidsFileKey]),
	}
}

func NewAirportFinder() *AirportFinder {
	return &AirportFinder{
		airportDB:   NewAirportDB(),
		frequencyDB: NewFrequencyDB(),
		runwayDB:    NewRunwayDB(),
		regionDB:    NewRegionDB(),
		countryDB:   NewCountryDB(),
		navaidDB:    NewNavaidDB(),
	}
}

func Clear(af *AirportFinder) {
	af.airportDB.Clear()
	af.frequencyDB.Clear()
	af.runwayDB.Clear()
	af.regionDB.Clear()
	af.countryDB.Clear()
	af.navaidDB.Clear()
}

func (af *AirportFinder) Load(options *LoadOptions, airportFilter uint64) []error {
	errors := make([]error, 0)
	if options == nil {
		return append(errors, fmt.Errorf("unable not load anything since options are null"))
	}
	if options.AirportsFilename == "" {
		return append(errors, fmt.Errorf("cannot load airports: invalid filename"))
	}

	if err := af.airportDB.Parse(options.AirportsFilename, airportFilter, true); err != nil {
		errors = append(errors, err)
	}
	if options.FrequenciesFilename != "" {
		if err := af.frequencyDB.Parse(options.FrequenciesFilename, true); err != nil {
			errors = append(errors, err)
		}
	}
	if options.RunwaysFilename != "" {
		if err := af.runwayDB.Parse(options.RunwaysFilename, true); err != nil {
			errors = append(errors, err)
		}
	}
	if options.RegionsFilename != "" {
		if err := af.regionDB.Parse(options.RegionsFilename, true); err != nil {
			errors = append(errors, err)
		}
	}
	if options.CountriesFilename != "" {
		if err := af.countryDB.Parse(options.CountriesFilename, true); err != nil {
			errors = append(errors, err)
		}
	}
	if options.NavaidsFilename != "" {
		if err := af.navaidDB.Parse(options.NavaidsFilename, true); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (af *AirportFinder) FindAirportByType(airportTypeFilter uint64) []*Airport {
	airportsByType := af.airportDB.FindByAirportType(airportTypeFilter)
	airports := make([]*Airport, 0, len(airportsByType))
	for _, airport := range airportsByType {
		airports = append(airports, af.makeAirport(airport))
	}
	return airports
}

func (af *AirportFinder) FindAirportByICAOCode(icaoCode string) *Airport {
	airport := af.airportDB.FindByICAOCode(icaoCode)
	return af.makeAirport(airport)
}

func (af *AirportFinder) FindAirportByIATACode(iataCode string) *Airport {
	airport := af.airportDB.FindByIATACode(iataCode)
	return af.makeAirport(airport)
}

func (af *AirportFinder) FindAirportsByRegion(isoRegion string, airportTypeFilter uint64) []*Airport {
	airportsByRegion := af.airportDB.FindByRegion(isoRegion, airportTypeFilter)
	airports := make([]*Airport, 0, len(airportsByRegion))
	for _, airport := range airportsByRegion {
		airports = append(airports, af.makeAirport(airport))
	}
	return airports
}

func (af *AirportFinder) FindAirportsByCountry(isoCountry string, airportTypeFilter uint64) []*Airport {
	airportsByCountry := af.airportDB.FindByCountry(isoCountry, airportTypeFilter)
	airports := make([]*Airport, 0, len(airportsByCountry))
	for _, airport := range airportsByCountry {
		airports = append(airports, af.makeAirport(airport))
	}
	return airports
}

func (af *AirportFinder) FindNearestAirport(latitudeDeg, longitudeDeg, radiusMeters float64, airportTypeFilter uint64) *Airport {
	nearestAirport := af.airportDB.FindNearestAirport(latitudeDeg, longitudeDeg, radiusMeters, airportTypeFilter)
	return af.makeAirport(nearestAirport)
}

func (af *AirportFinder) FindNearestAirports(latitudeDeg, longitudeDeg, radiusMeters float64, maxResults int, airportTypeFilter uint64) []*Airport {
	nearestAirports := af.airportDB.FindNearestAirports(latitudeDeg, longitudeDeg, radiusMeters, maxResults, airportTypeFilter)
	airports := make([]*Airport, 0, len(nearestAirports))
	for _, airport := range nearestAirports {
		airports = append(airports, af.makeAirport(airport))
	}
	return airports
}

func (af *AirportFinder) GetAllAirports(airportTypeFilter uint64) []*Airport {
	filteredAirports := af.airportDB.GetAllByFilter(airportTypeFilter)
	airports := make([]*Airport, 0, len(filteredAirports))
	for _, airport := range filteredAirports {
		airports = append(airports, af.makeAirport(airport))
	}
	return airports
}

func (af *AirportFinder) FindNearestNavaids(latitudeDeg, longitudeDeg, radiusMeters float64, maxResults int) []*Navaid {
	nearestNavaids := af.navaidDB.FindNearestNavaids(latitudeDeg, longitudeDeg, radiusMeters, maxResults)
	navaids := make([]*Navaid, 0, len(nearestNavaids))
	for _, navaid := range nearestNavaids {
		navaids = append(navaids, NewNavaid(navaid))
	}
	return navaids
}

func (af *AirportFinder) FindNavaidsByAirportICAOCode(icaoCode string) []*Navaid {
	associatedNavaids := af.navaidDB.FindByAirportICAOCode(icaoCode)
	navaids := make([]*Navaid, 0, len(associatedNavaids))
	for _, navaid := range associatedNavaids {
		navaids = append(navaids, NewNavaid(navaid))
	}
	return navaids
}

func (af *AirportFinder) GetAllNavaids() []*Navaid {
	navaids := make([]*Navaid, 0, len(af.navaidDB.Navaids))
	if af.navaidDB != nil {
		for _, navaid := range af.navaidDB.Navaids {
			navaids = append(navaids, NewNavaid(navaid))
		}
	}
	return navaids
}

func (af *AirportFinder) makeAirport(airport *AirportData) *Airport {
	if airport == nil {
		return nil
	}
	frequencies := af.frequencyDB.FindByAirportID(airport.ID)
	runways := af.runwayDB.FindByAirportID(airport.ID)
	region := af.regionDB.FindByISOCode(airport.ISORegion)
	country := af.countryDB.FindByISOCode(airport.ISOCountry)
	navaids := af.navaidDB.FindByAirportICAOCode(airport.ICAOCode)
	return NewAirport(airport, region, country, frequencies, runways, navaids)
}
