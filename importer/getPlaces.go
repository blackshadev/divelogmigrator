package importer

import (
	"database/sql"

	"littledev.nl/divelogimporter/models"
)

const SelectPlacesQuery = `
SELECT 
     id
   , name
  FROM places p
 WHERE EXISTS (SELECT 1 FROM dives d WHERE p.id = d.place_id AND d.user_id = $1)
 ORDER BY id ASC
`

func GetPlaces(sql *sql.DB, userId models.UserId) []models.Place {
	var places []models.Place
	rows, err := sql.Query(SelectPlacesQuery, int(userId))

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var place models.Place

		if err := rows.Scan(
			&place.Id,
			&place.Name,
		); err != nil {
			panic(err)
		}

		places = append(places, place)
	}

	return places
}
