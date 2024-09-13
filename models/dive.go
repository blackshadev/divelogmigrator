package models

import (
	"database/sql"
	"fmt"
	"time"
)

type DiveId uint

type Tank struct {
	Volume *uint8
	Begin  *float64
	End    *float64
	Oxygen *float64
}

type Dive struct {
	Id          DiveId
	Date        time.Time
	Divetime    Duration
	MaxDepth    Depth
	Tank        Tank
	ComputerId  Optional[ComputerId]
	PlaceId     Optional[PlaceId]
	Fingerprint Optional[uint32]
	Buddies     sql.NullString
	Tags        sql.NullString
}

type Depth float64

func (d Depth) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%.2f m", d)
	return []byte(s), nil
}

type Duration uint

func (d Duration) MarshalText() ([]byte, error) {
	mins := uint(d) / 60
	secs := uint(d) % 60

	return []byte(fmt.Sprintf("%02d:%02d min", mins, secs)), nil
}
