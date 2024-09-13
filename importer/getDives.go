package importer

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"

	"littledev.nl/divelogimporter/models"
)

const SelectDivesQuery = `
SELECT 
     id
   , date
   , d.divetime  
   , max_depth
   , (select string_agg(b.name, ', ') from buddy_dive bd join buddies b on bd.buddy_id  = b.id where bd.dive_id = d.id ) as buddies
   , (select string_agg(t.text , ', ') from dive_tag dt join tags t on dt.tag_id  = t.id where dt.dive_id = d.id ) as tags
   , (select dt2.oxygen from dive_tanks dt2 where dt2.dive_id = d.id limit 1) as tank_oxygen
   , (select dt2.volume from dive_tanks dt2 where dt2.dive_id = d.id limit 1) as tank_volume
   , (select dt2.pressure_begin from dive_tanks dt2 where dt2.dive_id = d.id limit 1) as tank_begin
   , (select dt2.pressure_end from dive_tanks dt2 where dt2.dive_id = d.id limit 1) as tank_end
   , d.computer_id
   , d.fingerprint
   , d.place_id
  FROM dives d
 WHERE user_id = $1
 ORDER BY date ASC
`

func GetDives(sql *sql.DB, userId models.UserId) []models.Dive {
	var dives []models.Dive
	rows, err := sql.Query(SelectDivesQuery, int(userId))

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var dive models.Dive
		var computerId *models.ComputerId
		var fingerprint *string
		var placeId *models.PlaceId

		if err := rows.Scan(
			&dive.Id,
			&dive.Date,
			&dive.Divetime,
			&dive.MaxDepth,
			&dive.Buddies,
			&dive.Tags,
			&dive.Tank.Oxygen,
			&dive.Tank.Volume,
			&dive.Tank.Begin,
			&dive.Tank.End,
			&computerId,
			&fingerprint,
			&placeId,
		); err != nil {
			panic(err)
		}

		dive.ComputerId = models.NewOptionalValue(computerId)
		dive.Fingerprint = parseDiveFingerprint(fingerprint)
		dive.PlaceId = models.NewOptionalValue(placeId)

		dives = append(dives, dive)
	}

	return dives
}

func parseDiveFingerprint(fingerprint *string) models.Optional[uint32] {
	if fingerprint == nil {
		return models.Optional[uint32]{
			IsFilled: false,
		}
	}

	bytes, err := base64.StdEncoding.DecodeString(*fingerprint)
	if err != nil {
		panic(err)
	}

	return models.Optional[uint32]{
		IsFilled: true,
		Value:    binary.BigEndian.Uint32(bytes),
	}
}
