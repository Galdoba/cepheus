package star

import (
	"errors"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"
)

/*
PostStellar: NS pNb BH D PSR
PreStellar : Protostar Nb

Nb - если выпадает для главной зыезды
pNb -если выпадает для вторичных звезд
Cluster - отрабатываем только главную звезду, если да то в системе должно быть 4+ звезд
Anomaly - делаем пометку, заполняем руками
*/

type Star struct {
	Type        string  `json:"type"`
	SubType     string  `json:"sub_type,omitempty"`
	Class       string  `json:"class,omitempty"`
	Mass        float64 `json:"mass,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Diameter    float64 `json:"diameter,omitempty"`
	Luminocity  float64 `json:"luminocity,omitempty"`
	Age         float64 `json:"age,omitempty"`
}

func NewPrecursor() *Precursor {
	return &Precursor{
		id: time.Now().UnixMicro(),
	}
}

type Precursor struct {
	id              int64
	starType        string
	starSubType     string
	starClass       string
	starDesignation stellar.StarDesignation
	starDead        bool
	starProtostar   bool
	starMass        float64
	starTemperature float64
	starDiameter    float64
	starLuminocity  float64
	starAge         float64
	starAnomaly     string
	starPeriod      float64
}

func (p *Precursor) Type() string {
	return p.starType
}

func (p *Precursor) SubType() string {
	return p.starSubType
}

func (p *Precursor) Class() string {
	return p.starClass
}

func (p *Precursor) Mass() float64 {
	return p.starMass
}

func (p *Precursor) Temperature() float64 {
	return p.starTemperature
}

func (p *Precursor) Diameter() float64 {
	return p.starDiameter
}

func (p *Precursor) Luminocity() float64 {
	return p.starLuminocity
}

func (p *Precursor) Age() float64 {
	return p.starAge
}

func (p *Precursor) Designation() stellar.StarDesignation {
	return p.starDesignation
}

func (p *Precursor) IsDead() bool {
	return p.starDead
}

func (p *Precursor) IsProtostar() bool {
	return p.starProtostar
}

func (p *Precursor) Anomaly() string {
	return p.starAnomaly
}

func (p *Precursor) Period() float64 {
	return p.starPeriod
}

func (p *Precursor) SetType(t string) {
	p.starType = t
}

func (p *Precursor) SetSubType(st string) {
	p.starSubType = st
}

func (p *Precursor) SetClass(c string) {
	p.starClass = c
}

func (p *Precursor) SetMass(m float64) {
	p.starMass = m
}

func (p *Precursor) SetTemperature(temp float64) {
	p.starTemperature = temp
}

func (p *Precursor) SetDiameter(d float64) {
	p.starDiameter = d
}

func (p *Precursor) SetLuminocity(l float64) {
	p.starLuminocity = l
}

func (p *Precursor) SetAge(a float64) {
	p.starAge = a
}

func (p *Precursor) SetDesignation(d stellar.StarDesignation) {
	p.starDesignation = d
}

func (p *Precursor) SetDead(dead bool) {
	p.starDead = dead
}

func (p *Precursor) SetProtostar(proto bool) {
	p.starProtostar = proto
}

func (p *Precursor) SetAnomaly(a string) {
	p.starAnomaly = a
}

func (p *Precursor) SetPeriod(per float64) {
	p.starPeriod = per
}

func Finalize(p *Precursor) (Star, error) {
	if p.starType == "" {
		return Star{}, errors.New("star type is empty: cannot finalize Star")
	}

	return Star{
		Type:        p.starType,
		SubType:     p.starSubType,
		Class:       p.starClass,
		Mass:        p.starMass,
		Temperature: p.starTemperature,
		Diameter:    p.starDiameter,
		Luminocity:  p.starLuminocity,
		Age:         p.starAge,
	}, nil
}

func (s Star) IsPostStellar() bool {
	switch s.Type {
	case "BH", "NS", "PSR", "D":
		return true
	}
	return false
}

func (s Star) IsPreStellar() bool {
	switch s.Type {
	case "Protostar", "BD", "L", "T", "Y":
		return true
	}
	return false
}
