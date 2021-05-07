package alphafoxtrot

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// https://ourairports.com/help/data-dictionary.html
// header: "id","airport_ref","airport_ident","type","description","frequency_mhz"
// e.g, 60768,3632,"KLAX","ATIS","ATIS",133.8

const (
	colFrequencyID = iota
	colFrequencyAirportRef
	colFrequencyAirportIdent
	colFrequencyType
	colFrequencyDescription
	colFrequencyMHZ
)

type FrequencyData struct {
	ID           uint64
	AirportID    uint64
	AirportIdent string
	Type         string
	Description  string
	FrequencyMHZ float64
}

type FrequencyDB struct {
	Frequencies map[uint64][]*FrequencyData
}

func NewFrequencyDB() *FrequencyDB {
	return &FrequencyDB{
		Frequencies: make(map[uint64][]*FrequencyData),
	}
}

func (db *FrequencyDB) Clear() {
	db.Frequencies = nil
}

func (db *FrequencyDB) Parse(file string, skipFirstLine bool) error {
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

		id, err := ParseUint(row[colFrequencyID])
		if err != nil {
			log.Println(line, err)
			continue
		}

		airportRef, err := ParseUint(row[colFrequencyAirportRef])
		if err != nil {
			log.Println(line, err)
			continue
		}

		airportIdent := row[colFrequencyAirportIdent]
		typ := row[colFrequencyType]
		desc := row[colFrequencyDescription]

		mhz, err := ParseFloat(row[colFrequencyMHZ])
		if err != nil {
			log.Println(line, err)
			continue
		}

		frequency := &FrequencyData{
			ID:           id,
			AirportID:    airportRef,
			AirportIdent: airportIdent,
			Type:         typ,
			Description:  desc,
			FrequencyMHZ: mhz,
		}
		db.Frequencies[airportRef] = append(db.Frequencies[airportRef], frequency)
	}
}

func (db *FrequencyDB) FindByAirportID(airportID uint64) []*FrequencyData {
	frequencies := make([]*FrequencyData, 0)
	list, ok := db.Frequencies[airportID]
	if ok {
		frequencies = append(frequencies, list...)
	}
	return frequencies
}
