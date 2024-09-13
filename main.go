package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/lib/pq"
	"littledev.nl/divelogimporter/context"
	"littledev.nl/divelogimporter/importer"
	"littledev.nl/divelogimporter/writer/subsurface"
)

func main() {

	path := flag.String("output", "output.xml", "file to output to")
	flag.Parse()

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "littledivelog", "littledivelog", "littledivelog"))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	sampler := importer.LLSampler{
		Sql: db,
	}

	userId := importer.SelectUser(db)
	computers := importer.GetComputers(db, userId)
	dives := importer.GetDives(db, userId)
	places := importer.GetPlaces(db, userId)

	context := context.CreateContext(&sampler, userId, computers, places, dives)

	fmt.Printf("Found %d dives\n", len(dives))

	writer := subsurface.Writer{
		TargetPath: *path,
	}

	writer.Write(context)

	fmt.Printf("Written to to %s\n", writer.TargetPath)
}
