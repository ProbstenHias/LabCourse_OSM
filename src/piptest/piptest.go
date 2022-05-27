package piptest

import (
	"log"
	"math"
	"sync"
	"time"
)

func getNewNorth() []float64 {
	return []float64{90, 0}
}

func degToRad(angle float64) float64 {
	return angle * (math.Pi / 180)
}

func getAngleToNorth(poleToBe []float64, toTransform []float64) float64 {

	latP := degToRad(poleToBe[0])
	lonP := degToRad(poleToBe[1])

	latA := degToRad(toTransform[0])
	lonA := degToRad(toTransform[1])

	if latP == math.Pi/2 {
		return lonA
	} else {

		y := math.Sin(lonA-lonP) * math.Cos(latA)
		x := math.Sin(latA)*math.Cos(latP) - math.Cos(latA)*math.Sin(latP)*math.Cos(lonA-lonP)

		return math.Atan2(y, x)
	}
}

func eastWest(c float64, d float64) int {
	delta := d - c

	if delta > math.Pi {
		delta = delta - 2*math.Pi
	}

	if delta < -math.Pi {
		delta = delta + 2*math.Pi
	}

	if delta > 0 && delta != math.Pi {
		return -1 // d west of c
	} else if delta < 0 && delta != -math.Pi {
		return 1 //d east of c
	} else {
		return 0 //d north or south of c (collinear)
	}

}

func checkPointP(point []float64, polygon [][]float64, tranNodes []float64) int {
	// return 1 P same as X, 0 for P != X, 2 P on edge, 3 antipodal P and X
	xLat := getNewNorth()[0]
	xLon := getNewNorth()[1]
	pLat := point[0]
	pLon := point[1]
	var i int
	var vBlat, vBlon, tlonB float64
	var vAlat, vAlon, tlonA float64

	// maybe we dont need this
	if pLat == -xLat {
		delLon := degToRad(pLon) - degToRad(xLon)
		if delLon < -math.Pi {
			delLon = delLon + 2*math.Pi
		}
		if delLon > math.Pi {
			delLon = delLon - 2*math.Pi
		}
		if delLon == math.Pi || delLon == -math.Pi {
			log.Printf("P (%f,%f) is antipodal to X (%f,%f). Cannot determine location", pLat, pLon, xLat, xLon)

			return 3 // return 3 for antipodal
		}
	}

	iCross := 0 //count crossings

	if degToRad(pLat) == degToRad(xLat) && degToRad(pLon) == degToRad(xLon) {
		return 1 // X same location as P
	}

	tLonP := getAngleToNorth(getNewNorth(), []float64{pLat, pLon})

	for i = 0; i < len(polygon); i++ {

		vAlat = polygon[i][0]
		vAlon = polygon[i][1]
		tlonA = tranNodes[i]

		if i < len(polygon)-1 {
			vBlat = polygon[i+1][0]
			vBlon = polygon[i+1][1]
			tlonB = tranNodes[i+1]
		} else {
			vBlat = polygon[0][0]
			vBlon = polygon[0][1]
			tlonB = tranNodes[0]
		}

		isTrike := 0

		if tLonP == tlonA {
			isTrike = 1
		} else {

			ewAB := eastWest(tlonA, tlonB)
			ewAP := eastWest(tlonA, tLonP)
			ewPB := eastWest(tLonP, tlonB)
			if ewAP == ewAB && ewPB == ewAB {
				isTrike = 1
			}
		}

		if isTrike == 1 {
			if pLat == vAlat && pLon == vAlon {
				return 2 //P lies on vertex of S
			}

			tLonX := getAngleToNorth([]float64{vAlat, vAlon}, []float64{xLat, xLon})
			tLonB := getAngleToNorth([]float64{vAlat, vAlon}, []float64{vBlat, vBlon})
			tLonP := getAngleToNorth([]float64{vAlat, vAlon}, []float64{pLat, pLon})

			if tLonP == tLonB {
				return 2 //P lies on side of S
			} else {
				ewBX := eastWest(tLonB, tLonX)
				ewBP := eastWest(tLonB, tLonP)

				if ewBX == -ewBP {
					iCross = iCross + 1
				}
			}
		}
	}

	if iCross%2 == 0 {
		return 1 // even number of times so P is where X is.
	}

	return 0
}

func isInBox(boundingBox []float64, pLoc []float64) bool { // p_loc (lat,long)

	return !(pLoc[0] < boundingBox[0] || pLoc[0] > boundingBox[1] || pLoc[1] < boundingBox[2] || pLoc[1] > boundingBox[3])

}

func isPointInWater(wayNodes [][][]float64, tranNodes [][]float64, boundBox [][]float64, point []float64) bool {
	//x in water some points might be antipodal. Run again with different x in that case.
	// return 1 if point in water. 0 otherwise
	var toRet int8 = 1

	for i, polygon := range wayNodes {

		if isInBox(boundBox[i], point) {

			loc := checkPointP(point, polygon, tranNodes[i])

			if loc == 0 || loc == 2 { //treating edges as land
				toRet = 0
				break
			} else {
				continue //check next polygon to see if crossing or not
			}
		}
	}

	return toRet == 1
}

//
func transformNodes(wayNodes [][][]float64) [][]float64 {
	tranNodes := make([][]float64, len(wayNodes))

	for i, polygon := range wayNodes {
		for _, coordinates := range polygon {
			tranNodes[i] = append(tranNodes[i], getAngleToNorth(getNewNorth(), coordinates))
		}
	}

	return tranNodes
}

func createBoundingBoxes(wayNodes [][][]float64) [][]float64 {
	boundBox := make([][]float64, len(wayNodes))

	for i, polygon := range wayNodes {
		var minLat = polygon[0][0]
		var maxLat = polygon[0][0]
		var minLon = polygon[0][1]
		var maxLon = polygon[0][1]

		for _, coordinates := range polygon {
			if minLat > coordinates[0] {
				minLat = coordinates[0]
			}

			if maxLat < coordinates[0] {
				maxLat = coordinates[0]
			}

			if minLon > coordinates[1] {
				minLon = coordinates[1]
			}
			if maxLon < coordinates[1] {
				maxLon = coordinates[1]
			}
		}

		boundBox[i] = []float64{minLat, maxLat, minLon, maxLon}

	}
	return boundBox
}

// method to call when we want to do this
func TopLevel(wayNodes [][][]float64, spherePointsArr [][]float64) [][]float64 {
	var i int
	var correctPArray [][]float64
	// The point above is in water
	start11 := time.Now()
	boundBoxes := createBoundingBoxes(wayNodes)

	tranNodes := transformNodes(wayNodes)

	end11 := time.Now()
	duration11 := end11.Sub(start11)
	log.Printf("Preprocessing of PIP: %s\n", duration11)

	////////////// without goroutines (sequential) ////////////////////
	// for i = 0; i < len(get_p_array); i++ {
	// 	start := time.Now()
	// 	var flag bool = false

	// 	if get_p_loc(wayNodes, array_tran_way_nodes, bound_box, get_p_array[i], x_loc) == 1 {
	// 		flag = true
	// 		correct_p_array = append(correct_p_array, get_p_array[i])
	// 	}

	// 	end := time.Now()
	// 	duration := end.Sub(start)
	// 	fmt.Printf("Time to find where P[%d] is: %s.  ", i, duration)
	// 	if flag {
	// 		fmt.Printf("In Water. \n")
	// 	} else {
	// 		fmt.Printf("In Land. \n")
	// 	}
	// }
	////////////// without goroutines (sequential) ////////////////////
	start1 := time.Now()
	results := make(chan []float64, len(spherePointsArr)) //channel for water points from of goroutines are stored here
	countChan := make(chan bool, len(spherePointsArr))
	var wg sync.WaitGroup //wait group for goroutines

	wg.Add(1)
	go func(lenPArray int) {
		defer wg.Done()
		chanLen := len(countChan)
		chanLenNow := 0
		var counter = 0

		for {
			chanLenNow = len(countChan)
			if chanLenNow == lenPArray {
				log.Printf("%d/%d\n", lenPArray, lenPArray)
				break
			}
			if chanLen < chanLenNow {
				counter += chanLenNow - chanLen
				log.Printf("%d/%d\r", counter, lenPArray)
				chanLen = chanLenNow
			}
		}
	}(len(spherePointsArr))

	for i = 0; i < len(spherePointsArr); i++ {

		wg.Add(1) // add to wait group

		go func(wayNodes [][][]float64, tranNodes [][]float64, boundBox [][]float64, pPoint []float64) { // call goroutine
			defer wg.Done()
			if isPointInWater(wayNodes, tranNodes, boundBox, pPoint) {
				results <- pPoint
			}
			countChan <- true // used for the counter when one point is assigned
		}(wayNodes, tranNodes, boundBoxes, spherePointsArr[i])

	}

	wg.Wait()      //wait for all to finish
	close(results) // close channel

	for point := range results { //append results
		correctPArray = append(correctPArray, point)
	}

	end1 := time.Now()
	duration1 := end1.Sub(start1)
	log.Printf("All points locations found in: %s\n", duration1)

	return correctPArray

}
