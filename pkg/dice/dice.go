package dice

type dice struct {
	edges  int
	result int
}

type Mod interface {
	Mod() int
	ModKey() string
}
