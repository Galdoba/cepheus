package orbital

type Orbit struct {
	distance          float64 //Diameter of parent body or OrbitN
	high              float64
	low               float64
	width             float64
	eccentricity      float64
	mass              float64
	period            float64 //hours
	rotation          string  //retrograde?
	parentDesignation string  //позиционный код родителя
	parentType        string  //тип первичного тела
	bodyType          string  //то что вращается
	designation       string  //позиционный код
	body              any     //то что вращается
}

func (o *Orbit) Distance() float64 {
	return o.distance
}
