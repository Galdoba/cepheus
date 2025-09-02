package travellermap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Galdoba/cepheus/pkg/grid/coordinates"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/cube"
)

const (
	DATA_URL = "https://www.travellermap.com/data"
)

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get responce error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http returned status: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("response read error: %v", err)
	}
	return data, nil
}

func jumpMapUrl(abb, hex string, distance int) string {
	return fmt.Sprintf("%v/%v/%v/jump/%v", DATA_URL, abb, hex, distance)
}

func UpdateSectorData() (SectorList, error) {
	sectors := SectorList{}
	data, err := get(DATA_URL)
	if err != nil {
		return sectors, fmt.Errorf("failed to get sector data: %v", err)
	}
	if len(data) == 0 {
		return sectors, fmt.Errorf("no data received")
	}
	if err := json.Unmarshal(data, &sectors); err != nil {
		return sectors, fmt.Errorf("failed to unmarshal sectors data: %v", err)
	}
	return sectors, nil
}

// GetWorldData downloads json data of jump map from travellermap.com and unmarshal it to World structs.
//
// base api link: https://travellermap.com/api/jumpworlds?x=x&y=y&jump=radius (radius = [0-12])
func GetWorldData(coords coordinates.SpaceCoordinates, radius int) ([]WorldData, error) {
	x, y := coords.GlobalValues()
	url := fmt.Sprintf("https://travellermap.com/api/jumpworlds?x=%v&y=%v&jump=%v", x, y, radius)
	data, err := get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get world data: %v", err)
	}
	list := WorldList{}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal world data: %v", err)
	}
	return list.Worlds, nil
}

type Database struct {
	Worlds map[string]WorldData
}

func FullMapUpdate(path string) error {
	crdList := calibrationPoints(13)
	listLen := len(crdList)
	database := Database{}
	database.Worlds = make(map[string]WorldData)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	worldsAdded := 0
	for batchNum, crd := range crdList {
		fmt.Printf("update in progress: %v/%v [worlds: %v]\r", batchNum+1, listLen, worldsAdded)
		worldBatch, err := GetWorldData(crd, 12)
		if err != nil {
			fmt.Printf("request %v failed: %v\n", batchNum+1, err)
			continue
		}
		for _, world := range worldBatch {
			key := fmt.Sprintf("{%v,%v}", world.WorldX, world.WorldY)
			database.Worlds[key] = world
			worldsAdded++
		}
	}
	data, err := json.MarshalIndent(&database, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func generateBeltCenters(n int) [][3]int {
	if n == 0 {
		return [][3]int{{0, 0, 0}}
	}

	total := 2 * n
	var centers [][3]int

	for i := -total; i <= total; i++ {
		rem := total - abs(i)
		for j := -rem; j <= rem; j++ {
			k := -i - j
			if abs(i)+abs(j)+abs(k) == total {
				centers = append(centers, [3]int{24 * i, 24 * j, 24 * k})
			}
		}
	}

	return centers
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func listCalibrationPoints() []coordinates.SpaceCoordinates {
	crdList := []coordinates.SpaceCoordinates{}
	for i := 0; i < 23; i++ {
		for _, crd := range generateBeltCenters(i) {
			// fmt.Println(n, i, j, crd)
			coords := coordinates.NewSpaceCoordinates(crd[0], crd[1], crd[2])
			crdList = append(crdList, coords)
		}
	}
	return crdList
}

func calibrationPoints(n int) []coordinates.SpaceCoordinates {
	zero := cube.NewCube(0, 0, 0)
	points := make(map[cube.Cube]bool)
	points[zero] = true
	for i := 0; i <= n; i++ {
		nextPoints := []cube.Cube{}
		for point := range points {
			p1 := cube.Move(cube.Move(point, cube.DirectionNorth, 13), cube.DirectionNorthEast, 12)
			p2 := cube.Move(cube.Move(point, cube.DirectionNorthEast, 13), cube.DirectionSouthEast, 12)
			p3 := cube.Move(cube.Move(point, cube.DirectionSouthEast, 13), cube.DirectionSouth, 12)
			p4 := cube.Move(cube.Move(point, cube.DirectionSouth, 13), cube.DirectionSouthWest, 12)
			p5 := cube.Move(cube.Move(point, cube.DirectionSouthWest, 13), cube.DirectionNorthWest, 12)
			p6 := cube.Move(cube.Move(point, cube.DirectionNorthWest, 13), cube.DirectionNorth, 12)
			nextPoints = append(nextPoints, p1, p2, p3, p4, p5, p6)
		}
		for _, np := range nextPoints {
			points[np] = true
		}

	}
	list := []coordinates.SpaceCoordinates{}
	for point := range points {
		list = append(list, coordinates.NewSpaceCoordinates(point.Q, point.R, point.S))
	}

	fmt.Println(len(list))
	return list
}
