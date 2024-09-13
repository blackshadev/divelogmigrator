package importer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"littledev.nl/divelogimporter/models"
)

type Sampler interface {
	GetDiveSamples(diveId models.DiveId) []models.DiveSample
}

type LLSampler struct {
	Sql *sql.DB
}

type LLDivePressures []struct {
	Tank     uint8
	Pressure float64
}

type LLDiveEvents []struct {
	Type  string
	Value uint
	Flag  uint
}

type LLDiveSample struct {
	Time        uint
	Pressure    LLDivePressures
	Depth       *float64
	Temperature *float64
	RBT         *uint
	Events      LLDiveEvents
}

func (p *LLDivePressures) UnmarshalJSON(b []byte) error {
	if strings.ToLower(string(b)) == "null" {
		return nil
	}

	if b[0] == '{' {
		var oldPressure struct {
			Tank  uint8
			Value float64
		}
		if err := json.Unmarshal(b, &oldPressure); err != nil {
			return err
		}

		*p = append(*p, struct {
			Tank     uint8
			Pressure float64
		}{Tank: oldPressure.Tank, Pressure: oldPressure.Value})

		return nil

	} else if b[0] == '[' {
		var oldData []struct {
			Tank     uint8
			Pressure float64
		}
		if err := json.Unmarshal(b, &oldData); err != nil {
			return err
		}
		copy(*p, oldData)

		return nil
	}

	return fmt.Errorf("unable to unmarshal JSON into Pressures")
}

func (p *LLDiveEvents) UnmarshalJSON(b []byte) error {
	if strings.ToLower(string(b)) == "null" {
		*p = LLDiveEvents{}
		return nil
	}

	var newData []struct {
		Type  string
		Value uint
		Flag  uint
	}
	if err := json.Unmarshal(b, &newData); err == nil {
		copy(*p, newData)
		return nil
	}

	var oldEvents []struct {
		Name  string
		Value uint
		Flags uint
	}
	if err := json.Unmarshal(b, &oldEvents); err != nil {
		return err
	}

	for _, event := range oldEvents {
		*p = append(*p, struct {
			Type  string
			Value uint
			Flag  uint
		}{
			Type:  event.Name,
			Flag:  event.Flags,
			Value: event.Value,
		})
	}

	return nil
}

func (sampler *LLSampler) GetDiveSamples(diveId models.DiveId) []models.DiveSample {
	rows, err := sampler.Sql.Query("select samples from dives where id = $1", diveId)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return []models.DiveSample{}
	}

	var samplesJson sql.NullString
	if err := rows.Scan(&samplesJson); err != nil {
		panic(err)
	}
	if !samplesJson.Valid {
		return []models.DiveSample{}
	}

	var llSamples []LLDiveSample
	if err := json.Unmarshal([]byte(samplesJson.String), &llSamples); err != nil {
		panic(err)
	}

	var diveSamples []models.DiveSample
	for _, llSample := range llSamples {
		var pressures []models.TankPressure
		for _, pressure := range llSample.Pressure {
			pressures = append(pressures, models.TankPressure{
				Tank:     pressure.Tank,
				Pressure: pressure.Pressure,
			})
		}

		var events []models.Event
		for _, event := range llSample.Events {
			events = append(events, models.Event{
				Type:  event.Type,
				Flag:  event.Flag,
				Value: event.Value,
			})
		}

		diveSamples = append(diveSamples, models.DiveSample{
			Time:        models.Duration(llSample.Time),
			Depth:       llSample.Depth,
			RBT:         llSample.RBT,
			Temperature: llSample.Temperature,
			Pressure:    pressures,
			Events:      events,
		})
	}

	return diveSamples

}
