package airports

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// https://ourairports.com/help/data-dictionary.html
// header: "id","airport_ref","airport_ident","length_ft","width_ft","surface","lighted","closed","le_ident","le_latitude_deg","le_longitude_deg","le_elevation_ft","le_heading_degT","le_displaced_threshold_ft","he_ident","he_latitude_deg","he_longitude_deg","he_elevation_ft","he_heading_degT","he_displaced_threshold_ft"
// 240922,3632,"KLAX",12091,150,"CON",1,0,"07L",33.9358,-118.419,119,83,,"25R",33.9399,-118.38,94,263,957

const (
	colRunwayID = iota
	colRunwayAirportRef
	colRunwayAirportIdent
	colRunwayLengthFt
	colRunwayWidthFt
	colRunwaySurface
	colRunwayLighted
	colRunwayClosed
	colRunwayLowEndIdent
	colRunwayLowEndLatitudeDeg
	colRunwayLowEndLongitudeDeg
	colRunwayLowEndElevationFt
	colRunwayLowEndHeadingDegT
	colRunwayLowEndDisplacedThresholdFt
	colRunwayHighEndIdent
	colRunwayHighEndLatitudeDeg
	colRunwayHighEndLongitudeDeg
	colRunwayHighEndElevationFt
	colRunwayHighEndHeadingDegT
	colRunwayHighEndDisplacedThresholdFt
)

type RunwayData struct {
	ID                          uint64
	AirportID                   uint64
	AirportIdent                string
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

type RunwayDB struct {
	Runways map[uint64][]*RunwayData
}

func NewRunwayDB() *RunwayDB {
	return &RunwayDB{
		Runways: make(map[uint64][]*RunwayData),
	}
}

func (db *RunwayDB) Clear() {
	db.Runways = nil
}

func (db *RunwayDB) Parse(file string, skipFirstLine bool) error {
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

		id, err := ParseUint(row[colRunwayID])
		if err != nil {
			log.Println(line, err)
			continue
		}

		airportRef, err := ParseUint(row[colRunwayAirportRef])
		if err != nil {
			log.Println(line, err)
			continue
		}

		airportIdent := row[colRunwayAirportIdent]
		length, _ := ParseInt(row[colRunwayLengthFt])
		width, _ := ParseInt(row[colRunwayWidthFt])
		surface := row[colRunwaySurface]
		lighted := ParseBool(row[colRunwayLighted])
		closed := ParseBool(row[colRunwayClosed])

		leIdent := row[colRunwayLowEndIdent]
		leLatitude, _ := ParseFloat(row[colRunwayLowEndLatitudeDeg])
		leLongitude, _ := ParseFloat(row[colRunwayLowEndLatitudeDeg])
		leElevation, _ := ParseInt(row[colRunwayLowEndElevationFt])
		leHeading, _ := ParseFloat(row[colRunwayLowEndHeadingDegT])
		leDisplacedThreshold, _ := ParseInt(row[colRunwayLowEndDisplacedThresholdFt])

		heIdent := row[colRunwayHighEndIdent]
		heLatitude, _ := ParseFloat(row[colRunwayHighEndLatitudeDeg])
		heLongitude, _ := ParseFloat(row[colRunwayHighEndLatitudeDeg])
		heElevation, _ := ParseInt(row[colRunwayHighEndElevationFt])
		heHeading, _ := ParseFloat(row[colRunwayHighEndHeadingDegT])
		heDisplacedThreshold, _ := ParseInt(row[colRunwayHighEndDisplacedThresholdFt])

		runway := &RunwayData{
			ID:                          id,
			AirportID:                   airportRef,
			AirportIdent:                airportIdent,
			LengthFt:                    length,
			WidthFt:                     width,
			Surface:                     surface,
			Lighted:                     lighted,
			Closed:                      closed,
			LowEndIdent:                 leIdent,
			LowEndLatitudeDeg:           leLatitude,
			LowEndLongitudeDeg:          leLongitude,
			LowEndElevationFt:           leElevation,
			LowEndHeadingDegT:           leHeading,
			LowEndDisplacedThresholdFt:  leDisplacedThreshold,
			HighEndIdent:                heIdent,
			HighEndLatitudeDeg:          heLatitude,
			HighEndLongitudeDeg:         heLongitude,
			HighEndElevationFt:          heElevation,
			HighEndHeadingDegT:          heHeading,
			HighEndDisplacedThresholdFt: heDisplacedThreshold,
		}
		db.Runways[airportRef] = append(db.Runways[airportRef], runway)
	}
}

func (db *RunwayDB) FindByAirportID(airportID uint64) []*RunwayData {
	runways := make([]*RunwayData, 0)
	list, ok := db.Runways[airportID]
	if ok {
		runways = append(runways, list...)
	}
	return runways
}
