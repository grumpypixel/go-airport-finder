package airports

type Airport struct {
	ICAOCode         string
	Type             string
	Name             string
	LatitudeDeg      float64
	LongitudeDeg     float64
	ElevationFt      int64
	Continent        string
	Municipality     string
	ScheduledService bool
	GPSCode          string
	IATACode         string
	LocalCode        string
	HomeLink         string
	WikipediaLink    string
	Keywords         string
	Region           Region
	Country          Country
	Runways          []Runway
	Frequencies      []Frequency
	Navaids          []Navaid
}

type Frequency struct {
	Type         string
	Description  string
	FrequencyMHZ float64
}

type Runway struct {
	LengthFt                    int64
	WidthFt                     int64
	Surface                     string
	Lighted                     bool
	Closed                      bool
	LowEndIdent                 string
	LowEndLatitudeDeg           float64
	LowEndLongitudeDeg          float64
	LowEndElevationFt           int64
	LowEndHeadingDegT           float64
	LowEndDisplacedThresholdFt  int64
	HighEndIdent                string
	HighEndLatitudeDeg          float64
	HighEndLongitudeDeg         float64
	HighEndElevationFt          int64
	HighEndHeadingDegT          float64
	HighEndDisplacedThresholdFt int64
}

type Region struct {
	ISOCode       string
	LocalCode     string
	Name          string
	WikipediaLink string
	Keywords      string
}

type Country struct {
	ISOCode       string
	Name          string
	Continent     string
	WikipediaLink string
	Keywords      string
}

type Navaid struct {
	Ident                string
	Name                 string
	Type                 string
	FrequencyKHZ         uint64
	LatitudeDeg          float64
	LongitudeDeg         float64
	ElevationFt          int64
	ISOCountry           string
	DMEFrequencyKHZ      uint64
	DMEChannel           string
	DMELatitudeDeg       float64
	DMELongitudeDeg      float64
	DMEElevationFt       int64
	SlavedVariationDeg   float64
	MagneticVariationDeg float64
	UsageType            string
	Power                string
	AssociatedAirport    string
}

func NewAirport(airport *AirportData, region *RegionData, country *CountryData, frequencies []*FrequencyData, runways []*RunwayData, navaids []*NavaidData) *Airport {
	if airport == nil {
		return nil
	}
	aeroport := &Airport{
		ICAOCode:         airport.IATACode,
		Type:             airport.Type,
		Name:             airport.Name,
		LatitudeDeg:      airport.LatitudeDeg,
		LongitudeDeg:     airport.LongitudeDeg,
		ElevationFt:      airport.ElevationFt,
		Continent:        airport.Continent,
		Municipality:     airport.Municipality,
		ScheduledService: airport.ScheduledService,
		GPSCode:          airport.GPSCode,
		IATACode:         airport.IATACode,
		LocalCode:        airport.LocalCode,
		HomeLink:         airport.HomeLink,
		WikipediaLink:    airport.WikipediaLink,
		Keywords:         airport.Keywords,
		Frequencies:      make([]Frequency, 0, len(frequencies)),
		Runways:          make([]Runway, 0, len(runways)),
		Navaids:          make([]Navaid, 0, len(navaids)),
	}
	if region != nil {
		aeroport.Region = *NewRegion(region)
	}
	if country != nil {
		aeroport.Country = *NewCountry(country)
	}
	for _, frequency := range frequencies {
		aeroport.Frequencies = append(aeroport.Frequencies, *NewFrequency(frequency))
	}
	for _, runway := range runways {
		aeroport.Runways = append(aeroport.Runways, *NewRunway(runway))
	}
	for _, navaid := range navaids {
		aeroport.Navaids = append(aeroport.Navaids, *NewNavaid(navaid))
	}
	return aeroport
}

func NewRegion(region *RegionData) *Region {
	return &Region{
		ISOCode:       region.ISOCode,
		LocalCode:     region.LocalCode,
		Name:          region.Name,
		WikipediaLink: region.WikipediaLink,
		Keywords:      region.Keywords,
	}
}

func NewCountry(country *CountryData) *Country {
	return &Country{
		ISOCode:       country.ISOCode,
		Name:          country.Name,
		Continent:     country.Continent,
		WikipediaLink: country.WikipediaLink,
		Keywords:      country.Keywords,
	}
}
func NewFrequency(frequency *FrequencyData) *Frequency {
	return &Frequency{
		Type:         frequency.Type,
		Description:  frequency.Description,
		FrequencyMHZ: frequency.FrequencyMHZ,
	}
}

func NewRunway(runway *RunwayData) *Runway {
	return &Runway{
		LengthFt:                    runway.LengthFt,
		WidthFt:                     runway.WidthFt,
		Surface:                     runway.Surface,
		Lighted:                     runway.Lighted,
		Closed:                      runway.Closed,
		LowEndIdent:                 runway.LowEndIdent,
		LowEndLatitudeDeg:           runway.LowEndLatitudeDeg,
		LowEndLongitudeDeg:          runway.LowEndLongitudeDeg,
		LowEndElevationFt:           runway.LowEndElevationFt,
		LowEndHeadingDegT:           runway.LowEndHeadingDegT,
		LowEndDisplacedThresholdFt:  runway.LowEndDisplacedThresholdFt,
		HighEndIdent:                runway.HighEndIdent,
		HighEndLatitudeDeg:          runway.HighEndLatitudeDeg,
		HighEndLongitudeDeg:         runway.HighEndLongitudeDeg,
		HighEndElevationFt:          runway.HighEndElevationFt,
		HighEndHeadingDegT:          runway.HighEndHeadingDegT,
		HighEndDisplacedThresholdFt: runway.HighEndDisplacedThresholdFt,
	}
}
func NewNavaid(navaid *NavaidData) *Navaid {
	return &Navaid{
		Ident:                navaid.Ident,
		Name:                 navaid.Name,
		Type:                 navaid.Type,
		FrequencyKHZ:         navaid.FrequencyKHZ,
		LatitudeDeg:          navaid.LatitudeDeg,
		LongitudeDeg:         navaid.LongitudeDeg,
		ElevationFt:          navaid.ElevationFt,
		ISOCountry:           navaid.ISOCountry,
		DMEFrequencyKHZ:      navaid.DMEFrequencyKHZ,
		DMEChannel:           navaid.DMEChannel,
		DMELatitudeDeg:       navaid.DMELatitudeDeg,
		DMELongitudeDeg:      navaid.DMELongitudeDeg,
		DMEElevationFt:       navaid.DMEElevationFt,
		SlavedVariationDeg:   navaid.SlavedVariationDeg,
		MagneticVariationDeg: navaid.MagneticVariationDeg,
		UsageType:            navaid.UsageType,
		Power:                navaid.Power,
		AssociatedAirport:    navaid.AssociatedAirport,
	}
}
