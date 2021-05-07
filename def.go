package alphafoxtrot

import (
	"math"
	"strconv"
	"strings"
)

const (
	AirportTypeUnknown      uint64 = 0x00
	AirportTypeClosed       uint64 = 0x01
	AirportTypeHeliport     uint64 = 0x02
	AirportTypeSeaplaneBase uint64 = 0x04
	AirportTypeSmall        uint64 = 0x08
	AirportTypeMedium       uint64 = 0x10
	AirportTypeLarge        uint64 = 0x20
	AirportTypeAll          uint64 = AirportTypeClosed | AirportTypeHeliport | AirportTypeSeaplaneBase | AirportTypeSmall | AirportTypeMedium | AirportTypeLarge
	AirportTypeActive       uint64 = AirportTypeHeliport | AirportTypeSeaplaneBase | AirportTypeSmall | AirportTypeMedium | AirportTypeLarge
	AirportTypeRunways      uint64 = AirportTypeSmall | AirportTypeMedium | AirportTypeLarge
)

const (
	AirportsFileKey    = "airports"
	FrequenciesFileKey = "frequencies"
	RunwaysFileKey     = "runways"
	RegionsFileKey     = "regions"
	CountriesFileKey   = "countries"
	NavaidsFileKey     = "navaids"
	OurAirportsBaseURL = "https://ourairports.com/data/"
)

const (
	DegToRad    float64 = math.Pi / 180.0
	EarthRadius float64 = 6371.0 * 1000.0
)

var OurAirportsFiles = map[string]string{
	AirportsFileKey:    "airports.csv",
	FrequenciesFileKey: "airport-frequencies.csv",
	RunwaysFileKey:     "runways.csv",
	RegionsFileKey:     "regions.csv",
	CountriesFileKey:   "countries.csv",
	NavaidsFileKey:     "navaids.csv",
}

func KilometersToMeters(km float64) float64 {
	return km * 1000.0
}

func MilesToMeters(mi float64) float64 {
	return mi * 1609.34
}

func NauticalMilesToMeters(nm float64) float64 {
	return nm * 1852.0
}

func MetersToKilometers(m float64) float64 {
	return m * 0.001
}

func MetersToMiles(m float64) float64 {
	return m * 0.000621373
}

func MetersToNauticalMiles(m float64) float64 {
	return m * 0.000539957
}

// see https://stackoverflow.com/questions/43167417/calculate-distance-between-two-points-in-leaflet
// returns the distance between to coordinates in meters
func Distance(fromLatitudeDeg, fromLongitudeDeg, toLatitudeDeg, toLongitudeDeg float64) float64 {
	lat1 := fromLatitudeDeg * DegToRad
	lon1 := fromLongitudeDeg * DegToRad
	lat2 := toLatitudeDeg * DegToRad
	lon2 := toLongitudeDeg * DegToRad

	dtLat := lat2 - lat1
	dtLon := lon2 - lon1

	a := math.Pow(math.Sin(dtLat*0.5), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dtLon*0.5), 2)
	c := 2 * math.Asin(math.Sqrt(a))
	return c * EarthRadius
}

func ParseFloat(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return value, err
}

func ParseInt(str string) (int64, error) {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func ParseUint(str string) (uint64, error) {
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func ParseBool(str string) bool {
	str = strings.ToLower(str)
	switch str {
	case "true":
		fallthrough
	case "yes":
		fallthrough
	case "1":
		return true
	}
	return false
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
