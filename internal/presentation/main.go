package main

import (
	"fmt"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates/cube"
	"github.com/Galdoba/cepheus/internal/presentation/api"
)

func main() {
	urlAdresses := []string{}
	for _, coods := range calibrationPoints(13) {
		urlAdresses = append(urlAdresses, urlFromCoords(coods))
	}

	start := time.Now()

	// Вариант 1: Простая версия
	results, errors := api.GetDataSimple(urlAdresses...)

	// Вариант 2: С повторами и мониторингом
	// results, errors := GetData(urls...)

	// Вариант 3: С отслеживанием прогресса
	// progressCh, getResults := api.GetDataWithProgress(urlAdresses...)
	// for progress := range progressCh {
	// 	fmt.Printf("Обработано: %.1f%%\n", progress.Percent)
	// }
	// results, errors := getResults()

	elapsed := time.Since(start)

	fmt.Printf("Время выполнения: %v\n", elapsed)
	fmt.Printf("Успешно: %d, Ошибок: %d\n", len(results), len(errors))
	fmt.Printf("Среднее время на запрос: %v\n",
		elapsed/time.Duration(len(urlAdresses)))
	time.Sleep(10 * time.Second)
	fmt.Println(string(results[urlAdresses[0]]))
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
	return list
}

func urlFromCoords(coords coordinates.SpaceCoordinates) string {
	x, y := coords.GlobalValues()
	url := fmt.Sprintf("https://travellermap.com/api/jumpworlds?x=%v&y=%v&jump=%v", x, y, 12)
	return url
}
