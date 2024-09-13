package writer

import (
	"context"

	"littledev.nl/divelogimporter/models"
)

type DiveWriter interface {
	Write(context context.Context, dives []models.Dive)
}
