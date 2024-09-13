package importer

import (
	"database/sql"
	"fmt"
	"slices"
	"strconv"

	"littledev.nl/divelogimporter/models"
	"littledev.nl/divelogimporter/prompt"
)

func SelectUser(sql *sql.DB) models.UserId {
	rows, err := sql.Query("select id, name from users")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		ids = append(ids, id)
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	answer := prompt.Question("Which userId would you like to export?> ")
	userId, err := strconv.Atoi(answer)
	if err != nil {
		panic(err)
	}
	if !slices.Contains(ids, userId) {
		panic(fmt.Errorf("no user with userId '%d' available", userId))
	}

	return models.UserId(userId)
}
