package main

//Description:
//You are given an array coordinates, coordinates[i] = [x, y], where [x, y] represents the coordinate of a point.
//Check if these points make a straight line in the XY plane.

//use (x3-x2)*(y2-y1) == (y3-y2)*(x2-x1) && (x3-x1)*(y2-y1) == (y3-y1)*(x2-x1)
//to check if (x3,y3) is at the lane of ( (x1,y1),(x2,y2) ).
func checkStraightLine(coordinates [][]int) bool {
	for n := 2; n < len(coordinates); n++ {
		diffy12 := coordinates[1][1] - coordinates[0][1]
		diffx12 := coordinates[1][0] - coordinates[0][0]
		if !((coordinates[n][0]-coordinates[0][0])*diffy12 == (coordinates[n][1]-coordinates[0][1])*diffx12 &&
			(coordinates[n][0]-coordinates[1][0])*diffy12 == (coordinates[n][1]-coordinates[1][1])*diffx12) {
			return false
		}
	}
	return true
}
