package astrogation

import "github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"

func (as *Astrogation) TradePathExist(source, destination coordinates.Cube) bool {
	if coordinates.Distance(source, destination) > 4 {
		return false
	}
	j4CrdList := coordinates.Spiral(source, 4)

	midpoints := []coordinates.Cube{}
	for _, midPoint := range j4CrdList {
		if coordinates.Distance(midPoint, source) <= 2 && coordinates.Distance(midPoint, destination) <= 2 {
			midpoints = append(midpoints, midPoint)
		}
	}

	for _, midPoint := range midpoints {
		if _, err := as.basic.Read(midPoint.ToGlobal().DatabaseKey()); err == nil {
			return true
		}
	}
	return false
}
