package travellermap

type SectorList struct {
	Sectors []SectorData
}

type SectorData struct {
	X            int
	Y            int
	Milieu       string
	Abbreviation string
	Tags         string
	Names        []SectorName
}

type SectorName struct {
	Text string
	Lang string
}
