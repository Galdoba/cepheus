package worlds

type World struct {
	Port        Port
	Size        Size
	Atmosphere  Atmosphere
	Hudrosphere Hudrosphere
	Population  Population
	Government  Government
	LawLevel    LawLevel
	TechLevel   TechLevel
}

type Port struct {
	Code string
}
type Size struct {
	Value int
}
type Atmosphere struct {
	Value int
}
type Hudrosphere struct {
	Value int
}
type Population struct {
	Value int
}
type Government struct {
	Value int
}
type LawLevel struct {
	Value int
}
type TechLevel struct {
	Value int
}
