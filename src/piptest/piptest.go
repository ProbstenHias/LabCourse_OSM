package piptest

import (
	"OSM/src/spherePoints"
	"fmt"
	"math"
	"sync"
	"time"
)

func toRad(angle float64) float64 {
	return angle * (math.Pi / 180)
}

func getAngleToNorth(poleToBe []float64, toTransform []float64) float64 {

	latP := toRad(poleToBe[0])
	lonP := toRad(poleToBe[1])

	latA := toRad(toTransform[0])
	lonA := toRad(toTransform[1])

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

func checkPointP(pLat float64, pLon float64, xLat float64, xLon float64, lineNodes [][]float64, arrayTranNodes []float64) int {
	// return 1 P same as X, 0 for P != X, 2 P on edge, 3 antipodal P and X
	var i int
	var vBlat, vBlon, tlonB float64
	var vAlat, vAlon, tlonA float64
	if pLat == -xLat {
		dellon := toRad(pLon) - toRad(xLon)
		if dellon < -math.Pi {
			dellon = dellon + 2*math.Pi
		}
		if dellon > math.Pi {
			dellon = dellon - 2*math.Pi
		}
		if dellon == math.Pi || dellon == -math.Pi {
			fmt.Printf("P (%f,%f) is antipodal to X (%f,%f). Cannot determine location", pLat, pLon, xLat, xLon)

			return 3 // return 3 for antipodal
		}
	}

	iCross := 0 //count crossings

	if toRad(pLat) == toRad(xLat) && toRad(pLon) == toRad(xLon) {
		return 1 // X same location as P
	}

	tLonP := getAngleToNorth([]float64{xLat, xLon}, []float64{pLat, pLon})

	for i = 0; i < len(lineNodes); i++ {

		vAlat = lineNodes[i][1]
		vAlon = lineNodes[i][0]
		tlonA = arrayTranNodes[i]

		if i < len(lineNodes)-1 {
			vBlat = lineNodes[i+1][1]
			vBlon = lineNodes[i+1][0]
			tlonB = arrayTranNodes[i+1]
		} else {
			vBlat = lineNodes[0][1]
			vBlon = lineNodes[0][0]
			tlonB = arrayTranNodes[0]
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

func isInBox(boundingBox []float64, pLoc []float64) bool { // p_loc (long,lat)

	if pLoc[0] < boundingBox[0] || pLoc[0] > boundingBox[1] || pLoc[1] < boundingBox[2] || pLoc[1] > boundingBox[3] {

		return false
	}

	return true

}

func getPLoc(wayNodes map[int64][][]float64, arrayTranWayNodes map[int64][]float64, boundBox map[int64][]float64, pLoc []float64, x_loc []float64) int8 {
	//x in water some points might be antipodal. Run again with different x in that case.
	// return 1 if point in water. 0 otherwise
	var toRet int8 = 1

	for key, polyNodes := range wayNodes {

		if isInBox(boundBox[key], pLoc) {

			loc := checkPointP(pLoc[1], pLoc[0], x_loc[1], x_loc[0], polyNodes, arrayTranWayNodes[key])

			for {
				if loc == 3 { // Implemented but never reached because the antipodal to (0,90) is (0,-90) and no points there (hard limited during creation)
					recalcArrayTranWayNodes := transformNodesPoly(polyNodes, []float64{x_loc[0] - 20, x_loc[1]})   //recalculate transformed polygon
					loc = checkPointP(pLoc[1], pLoc[0], x_loc[1], x_loc[0]-20, polyNodes, recalcArrayTranWayNodes) //X antipodal to P move X 20 degrees west, still water.
				} else {
					break
				}
			}

			if loc == 0 || loc == 2 { //treating edges as land
				toRet = 0
				break
			} else {
				continue //check next polygon to see if crossing or not
			}
		}
	}

	return toRet
}

func transformNodesPoly(polyNodes [][]float64, xLoc []float64) []float64 {
	var tranNodes []float64
	var i int

	for i = 0; i < len(polyNodes); i++ {
		tranNodes = append(tranNodes, getAngleToNorth([]float64{xLoc[1], xLoc[0]}, []float64{polyNodes[i][1], polyNodes[i][0]}))
	}
	return tranNodes
}

func transformNodes(nodes map[int64][]float64, xLoc []float64) map[int64]float64 {
	tranNodes := make(map[int64]float64)

	for key, node := range nodes {

		tranNodes[key] = getAngleToNorth([]float64{xLoc[1], xLoc[0]}, []float64{node[1], node[0]})
	}

	return tranNodes
}

func getBoundingBox(nodes map[int64][]float64, ways map[int64][]int64) map[int64][]float64 {
	boundBox := make(map[int64][]float64)

	for key, NodeIDs := range ways {
		var minLat = nodes[NodeIDs[0]][1]
		var maxLat = nodes[NodeIDs[0]][1]
		var minLon = nodes[NodeIDs[0]][0]
		var maxLon = nodes[NodeIDs[0]][0]

		for _, nodeId := range NodeIDs {
			if minLat > nodes[nodeId][1] {
				minLat = nodes[nodeId][1]
			}

			if maxLat < nodes[nodeId][1] {
				maxLat = nodes[nodeId][1]
			}

			if minLon > nodes[nodeId][0] {
				minLon = nodes[nodeId][0]
			}
			if maxLon < nodes[nodeId][0] {
				maxLon = nodes[nodeId][0]
			}
		}

		boundBox[key] = []float64{minLon, maxLon, minLat, maxLat}
	}
	return boundBox
}

func TopLevel(nodes map[int64][]float64, ways map[int64][]int64, noOfPoints int64) [][]float64 {
	var i int
	var correctPArray [][]float64
	xLoc := []float64{0, 90} //choose initial point with known location (long,lat)
	// The point above is in water

	start1 := time.Now()

	getPArray, _ := spherePoints.GeneratePointsOnSphere(noOfPoints)

	start11 := time.Now()

	boundBox := getBoundingBox(nodes, ways)

	tranNodes := transformNodes(nodes, xLoc)

	//transform from way nodes to a single vector in a map containing everything (needed to access the next and previous nodes)
	wayNodes := make(map[int64][][]float64)
	arrayTranWayNodes := make(map[int64][]float64)
	for key, val := range ways {

		for _, nodeId := range val {
			wayNodes[key] = append(wayNodes[key], nodes[nodeId])
			arrayTranWayNodes[key] = append(arrayTranWayNodes[key], tranNodes[nodeId])
		}
	}

	end11 := time.Now()
	duration11 := end11.Sub(start11)
	fmt.Printf("Preprocessing of PIP: %s\n", duration11)

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

	results := make(chan []float64, len(getPArray)) //channel for water points from of goroutines are stored here
	countChan := make(chan bool, len(getPArray))
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
				fmt.Printf("%d/%d\n", lenPArray, lenPArray)
				break
			}
			if chanLen < chanLenNow {
				counter += chanLenNow - chanLen
				fmt.Printf("%d/%d\r", counter, lenPArray)
				chanLen = chanLenNow
			}
		}
	}(len(getPArray))

	for i = 0; i < len(getPArray); i++ {

		wg.Add(1) // add to wait group

		go func(wayNodes map[int64][][]float64, arrayTranWayNodes map[int64][]float64, boundBox map[int64][]float64, pPoint []float64, xLoc []float64) { // call goroutine
			defer wg.Done()
			if getPLoc(wayNodes, arrayTranWayNodes, boundBox, pPoint, xLoc) == 1 {
				results <- pPoint
			}
			countChan <- true // used for the counter when one point is assigned
		}(wayNodes, arrayTranWayNodes, boundBox, getPArray[i], xLoc)

	}

	wg.Wait()      //wait for all to finish
	close(results) // close channel

	for point := range results { //append results
		correctPArray = append(correctPArray, point)
	}

	end1 := time.Now()
	duration1 := end1.Sub(start1)
	fmt.Printf("All points locations found in: %s\n", duration1)

	return correctPArray

}
