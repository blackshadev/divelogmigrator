package subsurface

import (
	"encoding/xml"

	"littledev.nl/divelogimporter/models"
)

type SSDivelog struct {
	XMLName  xml.Name   `xml:"divelog"`
	Program  string     `xml:"program,attr"`
	Version  uint8      `xml:"version,attr"`
	Settings SSSettings `xml:"settings"`
	Sites    []SSSite   `xml:"divesites>site"`
	Dives    []SSDive   `xml:"dives>dive"`
}

type SSSettings struct {
	Fingerprints []SSDeviceFingerprints
}
type SSDeviceFingerprints struct {
	XMLName  xml.Name `xml:"fingerprint"`
	Model    string   `xml:"model,attr"`
	Serial   string   `xml:"serial,attr"`
	DeviceId string   `xml:"deviceid,attr"`
	DiveId   string   `xml:"diveid,attr"`
	Data     string   `xml:"data,attr"`
}

type SSSite struct {
	Uuid string `xml:"uuid,attr"`
	Name string `xml:"name,attr"`
}

type SSExtraData struct {
	XMLName xml.Name `xml:"extradata"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

type SSDiveComputer struct {
	Model     string `xml:"model,attr"`
	DeviceId  string `xml:"deviceid,attr"`
	DiveId    string `xml:"diveid,attr"`
	ExtraData []SSExtraData
	Events    []SSEvent
	Samples   []SSSample
}

type SSDive struct {
	Number       uint            `xml:"number,attr"`
	Date         string          `xml:"date,attr"`
	Time         string          `xml:"time,attr"`
	Depth        SSDiveDepth     `xml:"depth"`
	Duration     models.Duration `xml:"duration,attr"`
	Buddy        string          `xml:"buddy,omitempty"`
	Tags         string          `xml:"tags,attr,omitempty"`
	Cylinder     SSCylinder      `xml:"cylinder"`
	DiveComputer *SSDiveComputer `xml:"divecomputer,omitempty"`
	DiveSiteId   string          `xml:"divesiteid,attr,omitempty"`
}

type SSDiveDepth struct {
	Max models.Depth `xml:"max,attr"`
}

type SSCylinder struct {
	Size            string `xml:"size,attr,omitempty"`
	WorkingPressure string `xml:"workpressure,attr,omitempty"`
	Description     string `xml:"description,attr,omitempty"`
	O2              string `xml:"o2,attr,omitempty"`
	Start           string `xml:"start,attr,omitempty"`
	End             string `xml:"end,attr,omitempty"`
}

type SSSample struct {
	XMLName     xml.Name        `xml:"sample"`
	Time        models.Duration `xml:"time,attr"`
	Depth       string          `xml:"depth,attr,omitempty"`
	Temperature string          `xml:"temp,attr,omitempty"`
	Pressure    string          `xml:"pressure,attr,omitempty"`
}

type SSEvent struct {
	XMLName xml.Name        `xml:"event"`
	Time    models.Duration `xml:"time,attr"`
	Name    string          `xml:"name,attr"`
	Type    uint            `xml:"type,attr,omitempty"`
	Flag    uint            `xml:"flag,attr,omitempty"`
	Value   string          `xml:"value,attr,omitempty"`
}
