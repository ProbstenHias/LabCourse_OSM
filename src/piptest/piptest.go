package piptest

import (
	"OSM/src/spherePoints"
	"fmt"
	"math"
	"sync"
	"time"
)

func to_rad(angle float64) float64 {
	return angle * (math.Pi / 180)
}

func get_angle_to_noth(pole_to_be []float64, to_transform []float64) float64 {

	lat_p := to_rad(pole_to_be[0])
	lon_p := to_rad(pole_to_be[1])

	lat_a := to_rad(to_transform[0])
	lon_a := to_rad(to_transform[1])

	if lat_p == math.Pi/2 {
		return lon_a
	} else {

		y := (math.Sin(lon_a-lon_p) * math.Cos(lat_a))
		x := (math.Sin(lat_a)*math.Cos(lat_p) - math.Cos(lat_a)*math.Sin(lat_p)*math.Cos(lon_a-lon_p))

		return math.Atan2(y, x)
	}
}

func east_west(c float64, d float64) int {
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

func check_point_p(p_lat float64, p_lon float64, x_lat float64, x_lon float64, lineNodes [][]float64, array_tran_nodes []float64) int {
	// return 1 P same as X, 0 for P != X, 2 P on edge, 3 antipodal P and X
	var i int
	var vBlat, vBlon, tlonB float64
	var vAlat, vAlon, tlonA float64
	if p_lat == -x_lat {
		dellon := to_rad(p_lon) - to_rad(x_lon)
		if dellon < -math.Pi {
			dellon = dellon + 2*math.Pi
		}
		if dellon > math.Pi {
			dellon = dellon - 2*math.Pi
		}
		if dellon == math.Pi || dellon == -math.Pi {
			fmt.Printf("P (%f,%f) is antipodal to X (%f,%f). Cannot determine location", p_lat, p_lon, x_lat, x_lon)

			return 3 // return 3 for antipodal
		}
	}

	icross := 0 //count crossings

	if to_rad(p_lat) == to_rad(x_lat) && to_rad(p_lon) == to_rad(x_lon) {
		return 1 // X same location as P
	}

	tlonP := get_angle_to_noth([]float64{x_lat, x_lon}, []float64{p_lat, p_lon})

	for i = 0; i < len(lineNodes); i++ {

		vAlat = lineNodes[i][1]
		vAlon = lineNodes[i][0]
		tlonA = array_tran_nodes[i]

		if i < len(lineNodes)-1 {
			vBlat = lineNodes[i+1][1]
			vBlon = lineNodes[i+1][0]
			tlonB = array_tran_nodes[i+1]
		} else {
			vBlat = lineNodes[0][1]
			vBlon = lineNodes[0][0]
			tlonB = array_tran_nodes[0]
		}

		istrike := 0

		if tlonP == tlonA {
			istrike = 1
		} else {

			ewAB := east_west(tlonA, tlonB)
			ewAP := east_west(tlonA, tlonP)
			ewPB := east_west(tlonP, tlonB)
			if ewAP == ewAB && ewPB == ewAB {
				istrike = 1
			}
		}

		if istrike == 1 {
			if p_lat == vAlat && p_lon == vAlon {
				return 2 //P lies on vertex of S
			}

			tlon_X := get_angle_to_noth([]float64{vAlat, vAlon}, []float64{x_lat, x_lon})
			tlon_B := get_angle_to_noth([]float64{vAlat, vAlon}, []float64{vBlat, vBlon})
			tlon_P := get_angle_to_noth([]float64{vAlat, vAlon}, []float64{p_lat, p_lon})

			if tlon_P == tlon_B {
				return 2 //P lies on side of S
			} else {
				ewBX := east_west(tlon_B, tlon_X)
				ewBP := east_west(tlon_B, tlon_P)

				if ewBX == -ewBP {
					icross = icross + 1
				}
			}
		}
	}

	if icross%2 == 0 {
		return 1 // even number of times so P is where X is.
	}

	return 0
}

func is_in_box(bounding_box []float64, p_loc []float64) bool { // p_loc (long,lat)

	if p_loc[0] < bounding_box[0] || p_loc[0] > bounding_box[1] || p_loc[1] < bounding_box[2] || p_loc[1] > bounding_box[3] {

		return false
	}

	return true

}

func get_p_loc(wayNodes map[int64][][]float64, array_tran_way_nodes map[int64][]float64, bound_box map[int64][]float64, p_loc []float64, x_loc []float64) int8 {
	//x in water some points might be antipodal. Run again with different x in that case.
	// return 1 if point in water. 0 otherwise
	var to_ret int8 = 1

	for key, polyNodes := range wayNodes {

		if is_in_box(bound_box[key], p_loc) {

			loc := check_point_p(p_loc[1], p_loc[0], x_loc[1], x_loc[0], polyNodes, array_tran_way_nodes[key])

			for {
				if loc == 3 { // Implemented but never reached because the antipodal to (0,90) is (0,-90) and no points there (hard limited during creation)
					recalc_array_tran_way_nodes := transform_nodes_poly(polyNodes, []float64{x_loc[0] - 20, x_loc[1]})     //recalculate transformed polygon
					loc = check_point_p(p_loc[1], p_loc[0], x_loc[1], x_loc[0]-20, polyNodes, recalc_array_tran_way_nodes) //X antipodal to P move X 20 degrees west, still water.
				} else {
					break
				}
			}

			if loc == 0 || loc == 2 { //treating edges as land
				to_ret = 0
				break
			} else {
				continue //check next polygon to see if crossing or not
			}
		}
	}

	return to_ret
}

func transform_nodes_poly(polyNodes [][]float64, x_loc []float64) []float64 {
	var tran_nodes []float64
	var i int

	for i = 0; i < len(polyNodes); i++ {
		tran_nodes = append(tran_nodes, get_angle_to_noth([]float64{x_loc[1], x_loc[0]}, []float64{polyNodes[i][1], polyNodes[i][0]}))
	}
	return tran_nodes
}

func transform_nodes(nodes map[int64][]float64, x_loc []float64) map[int64]float64 {
	tran_nodes := make(map[int64]float64)

	for key, node := range nodes {

		tran_nodes[key] = get_angle_to_noth([]float64{x_loc[1], x_loc[0]}, []float64{node[1], node[0]})
	}

	return tran_nodes
}

func get_bounding_box(nodes map[int64][]float64, ways map[int64][]int64) map[int64][]float64 {
	bound_box := make(map[int64][]float64)

	for key, NodeIDs := range ways {
		var minlat float64 = nodes[NodeIDs[0]][1]
		var maxlat float64 = nodes[NodeIDs[0]][1]
		var minlon float64 = nodes[NodeIDs[0]][0]
		var maxlon float64 = nodes[NodeIDs[0]][0]

		for _, nodeid := range NodeIDs {
			if minlat > nodes[nodeid][1] {
				minlat = nodes[nodeid][1]
			}

			if maxlat < nodes[nodeid][1] {
				maxlat = nodes[nodeid][1]
			}

			if minlon > nodes[nodeid][0] {
				minlon = nodes[nodeid][0]
			}
			if maxlon < nodes[nodeid][0] {
				maxlon = nodes[nodeid][0]
			}
		}

		bound_box[key] = []float64{minlon, maxlon, minlat, maxlat}
	}
	return bound_box
}

func Top_level(nodes map[int64][]float64, ways map[int64][]int64, no_of_points int64) [][]float64 {
	var i int
	var correct_p_array [][]float64
	x_loc := []float64{0, 90} //choose inital point with known location (long,lat)
	// The point above is in water

	start1 := time.Now()

	get_p_array := spherePoints.GeneratePointsOnSphere(no_of_points)

	start11 := time.Now()

	bound_box := get_bounding_box(nodes, ways)

	tran_nodes := transform_nodes(nodes, x_loc)

	//transform from way nodes to a single vector in a map containing everything (needed to access the next and previour nodes)
	wayNodes := make(map[int64][][]float64)
	array_tran_way_nodes := make(map[int64][]float64)
	for key, val := range ways {

		for _, nodeId := range val {
			wayNodes[key] = append(wayNodes[key], nodes[nodeId])
			array_tran_way_nodes[key] = append(array_tran_way_nodes[key], tran_nodes[nodeId])
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

	results := make(chan []float64, len(get_p_array)) //channel for water points from of goroutines are stored here
	count_chan := make(chan bool, len(get_p_array))
	var wg sync.WaitGroup //wait group for goroutines

	wg.Add(1)
	go func(len_p_array int) {
		defer wg.Done()
		chan_len := len(count_chan)
		chan_len_now := 0
		var counter int = 0

		for {
			chan_len_now = len(count_chan)
			if chan_len_now == len_p_array {
				fmt.Printf("%d/%d\n", len_p_array, len_p_array)
				break
			}
			if chan_len < chan_len_now {
				counter += (chan_len_now - chan_len)
				fmt.Printf("%d/%d\r", counter, len_p_array)
				chan_len = chan_len_now
			}
		}
	}(len(get_p_array))

	for i = 0; i < len(get_p_array); i++ {

		wg.Add(1) // add to wait group

		go func(wayNodes map[int64][][]float64, array_tran_way_nodes map[int64][]float64, bound_box map[int64][]float64, p_point []float64, x_loc []float64) { // call goroutine
			defer wg.Done()
			if get_p_loc(wayNodes, array_tran_way_nodes, bound_box, p_point, x_loc) == 1 {
				results <- p_point
			}
			count_chan <- true // used for the counter when one point is assigned
		}(wayNodes, array_tran_way_nodes, bound_box, get_p_array[i], x_loc)

	}

	wg.Wait()      //wait for all to finsih
	close(results) // close channel

	for point := range results { //append results
		correct_p_array = append(correct_p_array, point)
	}

	end1 := time.Now()
	duration1 := end1.Sub(start1)
	fmt.Printf("All points locations found in: %s\n", duration1)

	return correct_p_array

}
