package orbit

type Orbit struct {
	//DesignationCode is code based on orbit relative to star system central object
	//It is created by Orbit Designation Rules.
	DesignationCode string
	Distance        float64
	Eccentricity    float64
	Retrograde      bool
	CurrentPosition float64
}
