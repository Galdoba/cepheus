package starsystem

import (
	"encoding/json"
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/orbit"
)

type StarSystem struct {
	SectorType      string
	Bodies          []SystemBody
	celestialObject string
	stars           []*Star
}

// prepend adds elem to the front of slice and returns a new slice.
func prepend[T any](slice []T, elem T) []T {
	return append([]T{elem}, slice...)
}

// insert inserts elem at the given index (0-based) and returns a new slice.
// Assumes index is valid (0 ≤ index ≤ len(slice)).
func insert[T any](index int, slice []T, elem T) []T {
	res := make([]T, 0, len(slice)+1)
	res = append(res, slice[:index]...)
	res = append(res, elem)
	res = append(res, slice[index:]...)
	return res
}

// replace returns a new slice where the element at index is replaced with elem.
// Assumes index is valid (0 ≤ index < len(slice)).
func replace[T any](index int, slice []T, elem T) []T {
	res := make([]T, len(slice))
	copy(res, slice)
	res[index] = elem
	return res
}

type SystemBody interface {
	json.Marshaler
	GetDesignation() string
	GetType() string
	GetOrbit() *orbit.Orbit
	GetChildren() []SystemBody
	SetChildren([]SystemBody)
}

func unmarshalBody(raw json.RawMessage) (SystemBody, error) {
	var descr struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &descr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal descriptor: %w", err)
	}

	switch descr.Type {
	case "star":
		s := Star{}
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, fmt.Errorf("failed to unmarshal star: %w", err)
		}
		return &s, nil
	}
	return nil, fmt.Errorf("unimplemented descriptor %q", descr.Type)
}

func UnmarshalBodies(data []byte) ([]SystemBody, error) {
	var rawSlice []json.RawMessage
	if err := json.Unmarshal(data, &rawSlice); err != nil {
		return nil, err
	}

	bodies := make([]SystemBody, 0, len(rawSlice))
	for _, raw := range rawSlice {
		b, err := unmarshalBody(raw)
		if err != nil {
			return nil, err
		}
		bodies = append(bodies, b)
	}
	return bodies, nil
}

func (ss *StarSystem) StarsPresent() int {
	n := 0
	for _, body := range ss.Bodies {
		if body.GetType() == "star" {
			n++
		}
	}
	return n
}

func (ss *StarSystem) AddStar(starType string) {
	str := NewStar()
	switch starType {
	case star:
		str.StellarClass = starType
	case brownDwarf:
		str.StellarClass = "BD"
	case neutronStar:
		str.StellarClass = "NS"
	case blachHole:
		str.StellarClass = "BH"
	}
	ss.stars = append(ss.stars, str)
}
