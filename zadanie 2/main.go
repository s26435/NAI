package main

import (
	"fmt"
	"math"
)

func FuzzifyHours(hours float64) (float64, error) {
	if hours <= 1 {
		return 1.0, fmt.Errorf("too high value of hours, should be (0 - 12)")
	} else if hours >= 12 {
		return 0.0, fmt.Errorf("too low value of hours, should be (0 - 12)")
	}
	return 1 - (hours-1)/11.0, nil
}

func FuzzifyImportance(importance float64) (float64, error) {
	if importance <= 0 {
		return 0.0, fmt.Errorf("too high value of importacne, should be (0 - 5)")
	} else if importance > 5 {
		return 1.0, fmt.Errorf("too high value of importacne, should be (0 - 5)")
	}
	return importance / 5.0, nil
}

func FuzzifyDistance(distance float64) (float64, error) {
	if distance <= 0 {
		return 1.0, fmt.Errorf("too high value of importacne, should be (0 - 5)")
	} else if distance >= 100 {
		return 0.0, fmt.Errorf("too high value of importacne, should be (0 - 5)")
	}
	return 1 - distance/100.0, nil
}

func CalculateAttendanceScore(hours, importance, distance float64) float64 {
	fHours, err := FuzzifyHours(hours)
	if err != nil {
		fmt.Printf("%v\n", err)
		return -1.0
	}
	fImportance, err := FuzzifyImportance(importance)
	if err != nil {
		fmt.Printf("%v\n", err)
		return -1.0
	}
	fDistance, err := FuzzifyDistance(distance)
	if err != nil {
		fmt.Printf("%v\n", err)
		return -1.0
	}
	score1 := math.Min(math.Min(fHours, fImportance), fDistance) * 5.0
	score2 := math.Min(fHours, fImportance) * 4.0
	score3 := math.Min(fImportance, fDistance) * 3.0
	totalScore := (score1 + score2 + score3) / 3.0
	return totalScore
}

func main() {
	hours := 6.0
	importance := 5.0
	distance := 15.0

	score := CalculateAttendanceScore(hours, importance, distance)
	fmt.Printf("Attendance worthiness score: %.2f out of 5\n", score)
}
