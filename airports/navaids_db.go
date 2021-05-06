package airports

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"sort"
)

// https://ourairports.com/help/data-dictionary.html
// "id","filename","ident","name","type","frequency_khz","latitude_deg","longitude_deg","elevation_ft","iso_country","dme_frequency_khz","dme_channel","dme_latitude_deg","dme_longitude_deg","dme_elevation_ft","slaved_variation_deg","magnetic_variation_deg","usageType","power","associated_airport"
// e.g. 90184,"Los_Angeles_VORTAC_US","LAX","Los Angeles","VORTAC",113600,33.933101654052734,-118.43199920654297,182,"US",113600,"083X",33.9334,-118.432,180,15.001,13.076,"BOTH","HIGH","KLAX"

const (
	colNavaidID = iota
	colNavaidFilename
	colNavaidIdent
	colNavaidName
	colNavaidType
	colNavaidFrequencyKHZ
	colNavaidLatitudeDeg
	colNavaidLongitudeDeg
	colNavaidElevationFt
	colNavaidISOCountry
	colNavaidDMEFrequencyKHZ
	colNavaidDMEChannel
	colNavaidDMELatitudeDeg
	colNavaidDMELongitudeDeg
	colNavaidDMEElevationFt
	colNavaidSlavedVariationDeg
	colNavaidMagneticVariationDeg
	colNavaidUsageType
	colNavaidPower
	colNavaidAssociatedAirport
)

type NavaidData struct {
	ID                   uint64
	Filename             string
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

type NavaidDB struct {
	Navaids []*NavaidData
}

func NewNavaidDB() *NavaidDB {
	return &NavaidDB{
		Navaids: make([]*NavaidData, 0),
	}
}

func (db *NavaidDB) Clear() {
	db.Navaids = nil
}

func (db *NavaidDB) Parse(file string, skipFirstLine bool) error {
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

		id, err := ParseUint(row[colNavaidID])
		if err != nil {
			log.Println(line, err)
			continue
		}

		filename := row[colNavaidFilename]
		ident := row[colNavaidIdent]
		name := row[colNavaidName]
		typ := row[colNavaidType]
		frequency, _ := ParseUint(row[colNavaidFrequencyKHZ])
		latitude, _ := ParseFloat(row[colNavaidLatitudeDeg])
		longitude, _ := ParseFloat(row[colNavaidLongitudeDeg])
		elevation, _ := ParseInt(row[colNavaidElevationFt])
		country := row[colNavaidISOCountry]

		dmeFrequency, _ := ParseUint(row[colNavaidDMEFrequencyKHZ])
		dmeChannel := row[colNavaidDMEChannel]

		dmeLatitude, _ := ParseFloat(row[colNavaidDMELatitudeDeg])
		dmeLongitude, _ := ParseFloat(row[colNavaidDMELongitudeDeg])
		dmeElevation, _ := ParseInt(row[colNavaidDMEElevationFt])

		slavedVariation, _ := ParseFloat(row[colNavaidSlavedVariationDeg])
		magneticVariation, _ := ParseFloat(row[colNavaidMagneticVariationDeg])

		usageType := row[colNavaidUsageType]
		power := row[colNavaidPower]
		associatedAirport := row[colNavaidAssociatedAirport]

		navaid := &NavaidData{
			ID:                   id,
			Filename:             filename,
			Ident:                ident,
			Name:                 name,
			Type:                 typ,
			FrequencyKHZ:         frequency,
			LatitudeDeg:          latitude,
			LongitudeDeg:         longitude,
			ElevationFt:          elevation,
			ISOCountry:           country,
			DMEFrequencyKHZ:      dmeFrequency,
			DMEChannel:           dmeChannel,
			DMELatitudeDeg:       dmeLatitude,
			DMELongitudeDeg:      dmeLongitude,
			DMEElevationFt:       dmeElevation,
			SlavedVariationDeg:   slavedVariation,
			MagneticVariationDeg: magneticVariation,
			UsageType:            usageType,
			Power:                power,
			AssociatedAirport:    associatedAirport,
		}
		db.Navaids = append(db.Navaids, navaid)
	}
}

func (db *NavaidDB) FindByAirportICAOCode(icaoCode string) []*NavaidData {
	navaids := make([]*NavaidData, 0)
	for _, navaid := range db.Navaids {
		if navaid.AssociatedAirport == icaoCode {
			navaids = append(navaids, navaid)
		}
	}
	return navaids
}

func (db *NavaidDB) FindNearestNavaids(latitudeDeg, longitudeDeg, radiusMeters float64, maxResults int) []*NavaidData {
	if radiusMeters < 0 {
		radiusMeters = math.MaxFloat64
	}
	if maxResults < 0 {
		maxResults = math.MaxInt32
	}

	type Candidate struct {
		Navaid   *NavaidData
		Distance float64
	}

	candidates := make([]*Candidate, 0)
	for _, navaid := range db.Navaids {
		distance := Distance(latitudeDeg, longitudeDeg, navaid.LatitudeDeg, navaid.LongitudeDeg)
		if distance <= radiusMeters {
			candidates = append(candidates, &Candidate{navaid, distance})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Distance < candidates[j].Distance
	})

	navaids := make([]*NavaidData, 0)
	count := MinInt(len(candidates), maxResults)
	for i := 0; i < count; i++ {
		navaids = append(navaids, candidates[i].Navaid)
	}
	return navaids
}
