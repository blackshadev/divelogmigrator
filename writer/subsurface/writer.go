package subsurface

import (
	"bufio"
	"crypto/sha1"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"os"

	"littledev.nl/divelogimporter/context"
	"littledev.nl/divelogimporter/models"
)

type Writer struct {
	TargetPath string
}

func formatDiveId(dive *models.Dive) string {
	if !dive.Fingerprint.IsFilled {
		return ""
	}

	return fmt.Sprintf("%x", sha1ifyUint(dive.Fingerprint.Value))
}

func formatDeviceId(computer *models.Computer) string {
	return fmt.Sprintf("%x", sha1ifyString(fmt.Sprintf("%d", computer.Serial)))
}

func getDiveByFingerprint(context context.Context, fingerprint uint32) *models.Dive {
	for _, dive := range context.Dives {
		if dive.Fingerprint.IsFilled && dive.Fingerprint.Value == fingerprint {
			return &dive
		}
	}
	return nil
}

func sha1ifyUint(i uint32) uint32 {
	var a = make([]byte, 4)
	binary.BigEndian.PutUint32(a, i)

	ahash := sha1.Sum(a)

	return binary.LittleEndian.Uint32(ahash[0:4])
}

func sha1ifyString(i string) uint32 {
	ahash := sha1.Sum([]byte(i))

	return binary.LittleEndian.Uint32(ahash[0:4])
}

func formatSiteId(place *models.Place) string {
	return fmt.Sprintf("%x", sha1ifyString(place.Name))
}

func translateSites(context context.Context) []SSSite {
	sites := []SSSite{}
	for _, place := range context.Places {
		sites = append(sites, SSSite{
			Uuid: formatSiteId(&place),
			Name: place.Name,
		})
	}

	return sites
}

func translateDeviceFingerprints(context context.Context) []SSDeviceFingerprints {
	var fingerprints []SSDeviceFingerprints

	for _, computer := range context.Computers {
		dive := getDiveByFingerprint(context, computer.Fingerprint)
		diveId := ""
		if dive != nil {
			diveId = formatDiveId(dive)
		}

		fingerprints = append(fingerprints, SSDeviceFingerprints{
			Model:    fmt.Sprintf("%x", sha1ifyString(computer.FullName())),
			Serial:   fmt.Sprintf("%x", computer.Serial),
			DeviceId: formatDeviceId(&computer),
			DiveId:   diveId,
			Data:     fmt.Sprintf("%x", computer.Fingerprint),
		})
	}

	return fingerprints
}

func translateExtraData(computer *models.Computer) []SSExtraData {
	return []SSExtraData{
		{
			Key:   "Serial",
			Value: fmt.Sprintf("%d", computer.Serial),
		},
	}
}

func translateDiveComputerData(context context.Context, dive models.Dive) *SSDiveComputer {
	samples := context.Sampler.GetDiveSamples(dive.Id)
	computer := context.GetComputer(dive.ComputerId)

	if !computer.IsFilled {
		return nil
	}

	return &SSDiveComputer{
		Model:     computer.Value.FullName(),
		DeviceId:  formatDeviceId(&computer.Value),
		DiveId:    formatDiveId(&dive),
		ExtraData: translateExtraData(&computer.Value),
		Events:    translateEvents(samples),
		Samples:   translateSamples(samples),
	}
}

func translateDives(context context.Context) []SSDive {
	dives := []SSDive{}
	for num, dive := range context.Dives {
		place := context.GetPlace(dive.PlaceId)

		diveSiteId := ""
		if place.IsFilled {
			diveSiteId = formatSiteId(&place.Value)
		}

		dives = append(dives, SSDive{
			Number:       uint(num + 1),
			Buddy:        dive.Buddies.String,
			Tags:         dive.Tags.String,
			Date:         dive.Date.Format("2006-01-02"),
			Time:         dive.Date.Format("15:04:05"),
			Duration:     dive.Divetime,
			Depth:        SSDiveDepth{Max: dive.MaxDepth},
			Cylinder:     translateTank(dive.Tank),
			DiveComputer: translateDiveComputerData(context, dive),
			DiveSiteId:   diveSiteId,
		})
	}
	return dives
}

func (w *Writer) Write(context context.Context) {
	sslog := &SSDivelog{
		Program: "subsurface",
		Version: 3,
		Settings: SSSettings{
			Fingerprints: translateDeviceFingerprints(context),
		},
		Sites: translateSites(context),
		Dives: translateDives(context),
	}

	out, err := xml.MarshalIndent(sslog, " ", "  ")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(w.TargetPath)
	if err != nil {
		panic(err)
	}
	io := bufio.NewWriter(f)
	io.Write(out)
	io.Flush()
}
