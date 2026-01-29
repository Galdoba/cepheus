package starsystem

import (
	"time"

	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/abstractions"
	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/orbital"
	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/star"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
)

type StarSystem struct {
	ID                string
	Name              string
	Designation       string
	CreatedAt         time.Time
	UpdatedAt         time.Time `json:"updated_at"`
	SurveyIndex       int
	ImportedData      t5ss.WorldData
	PrimaryStar       *star.Star
	StellarObjects    map[orbital.Orbit]abstractions.StellarObject
	WorldObjects      map[orbital.Orbit]abstractions.WorldObject
	PresenceConfirmed bool
}

type SystemBuilder struct {
	ImportedStellar string
}
