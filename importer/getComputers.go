package importer

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"

	"littledev.nl/divelogimporter/models"
)

const SelectComputerQuery = `
	SELECT id
		 , serial
		 , type
		 , vendor
		 , name
		 , last_fingerprint
	  FROM computers
	 WHERE user_id = $1 
`

func GetComputers(sql *sql.DB, userId models.UserId) []models.Computer {
	var computers []models.Computer
	rows, err := sql.Query(SelectComputerQuery, int(userId))

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		computer := models.Computer{}

		var lastFingerprint string
		rows.Scan(
			&computer.Id,
			&computer.Serial,
			&computer.Type,
			&computer.Vendor,
			&computer.Name,
			&lastFingerprint,
		)

		bytes, err := base64.StdEncoding.DecodeString(lastFingerprint)
		if err != nil {
			panic(err)
		}
		computer.Fingerprint = binary.BigEndian.Uint32(bytes)

		computers = append(computers, computer)
	}

	return computers
}
