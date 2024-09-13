package context

import (
	"littledev.nl/divelogimporter/importer"
	"littledev.nl/divelogimporter/models"
)

type Context struct {
	UserId    models.UserId
	Computers map[models.ComputerId]models.Computer
	Places    map[models.PlaceId]models.Place
	Dives     []models.Dive
	Sampler   importer.Sampler
}

func CreateContext(sampler importer.Sampler, userId models.UserId, computers []models.Computer, places []models.Place, dives []models.Dive) Context {
	computersById := make(map[models.ComputerId]models.Computer)
	placesById := make(map[models.PlaceId]models.Place)

	for _, computer := range computers {
		computersById[computer.Id] = computer
	}
	for _, place := range places {
		placesById[place.Id] = place
	}

	return Context{
		UserId:    userId,
		Computers: computersById,
		Dives:     dives,
		Sampler:   sampler,
		Places:    placesById,
	}
}

func (ctx *Context) GetPlace(id models.Optional[models.PlaceId]) models.Optional[models.Place] {
	if !id.IsFilled {
		return models.Optional[models.Place]{
			IsFilled: false,
		}
	}

	if place, ok := ctx.Places[id.Value]; ok {
		return models.Optional[models.Place]{
			IsFilled: true,
			Value:    place,
		}
	}

	return models.Optional[models.Place]{
		IsFilled: false,
	}
}

func (ctx *Context) GetComputer(id models.Optional[models.ComputerId]) models.Optional[models.Computer] {
	if !id.IsFilled {
		return models.Optional[models.Computer]{
			IsFilled: false,
		}
	}

	if computer, ok := ctx.Computers[id.Value]; ok {
		return models.Optional[models.Computer]{
			IsFilled: true,
			Value:    computer,
		}
	}

	return models.Optional[models.Computer]{
		IsFilled: false,
	}
}
